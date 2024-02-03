package sniper

import (
	"context"
	"errors"
	"fmt"
	"github.com/fidesy/me-sniper/internal/config"
	"github.com/fidesy/me-sniper/internal/pkg/crypto"
	"github.com/fidesy/me-sniper/internal/pkg/models"
	"github.com/gagliardetto/solana-go"
	rpc_ "github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/rpc"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	MEPublicKeyStr = "1BWutmTvYPwDtmw9abTkS4Ssr8no61spGAvW1X6NDix"
	MEPublicKey    = solana.MustPublicKeyFromBase58(MEPublicKeyStr)
)

type (
	Service struct {
		tokens              map[string]*models.Token
		solClient           *client.Client
		meClient            MagicEdenClient
		notificationService NotificationService
	}

	MagicEdenClient interface {
		GetFloor(ctx context.Context, symbol string) (float64, error)
		Buy(ctx context.Context, token *models.Token) (string, error)
	}

	NotificationService interface {
		SendNotification(ctx context.Context, action *models.Action)
	}
)

type Option func(s *Service)

func New(
	meClient MagicEdenClient,
	notificationService NotificationService,
	options ...Option) (*Service, error) {
	s := &Service{
		meClient:            meClient,
		notificationService: notificationService,
	}

	solanaEndpoint := config.Get(config.SolanaEndpoint).(string)
	if solanaEndpoint == "" {
		return nil, errors.New("solana endpoint config string is required")
	}

	solClient := client.NewClient(solanaEndpoint)
	// node health check
	if _, err := solClient.GetBalance(context.Background(), MEPublicKeyStr); err != nil {
		return nil, err
	}
	s.solClient = solClient

	tokens, err := models.LoadTokens()
	if err != nil {
		return nil, fmt.Errorf("models.LoadTokens: %w", err)
	}

	s.tokens = tokens

	for _, opt := range options {
		opt(s)
	}

	return s, nil
}

func (s *Service) Run(ctx context.Context) error {
	// For websocket connection public wss will be enough
	wsClient, err := ws.Connect(ctx, rpc_.MainNetBeta_WS)
	if err != nil {
		return fmt.Errorf("ws.Connect: %w", err)
	}

	sub, err := wsClient.LogsSubscribeMentions(MEPublicKey, "confirmed")
	if err != nil {
		return fmt.Errorf("wsClient.LogsSubscribeMentions: %w", err)
	}
	defer sub.Unsubscribe()

	go func() {
		for {
			got, err := sub.Recv()
			if err != nil {
				log.Printf("sub.Recv: %v", err)
				return
			}

			if got == nil {
				continue
			}

			go s.GetTransaction(ctx, got.Value.Signature.String())
		}
	}()

	<-ctx.Done()

	return nil
}

func (s *Service) GetTransaction(ctx context.Context, signature string) {
	var (
		transaction *client.GetTransactionResponse
		err         error
	)
	for transaction == nil {
		transaction, err = s.solClient.GetTransactionWithConfig(
			ctx,
			signature,
			rpc.GetTransactionConfig{Commitment: "confirmed"},
		)
		if err != nil || transaction == nil {
			time.Sleep(time.Millisecond * 500)
			continue
		}
	}

	token := s.parseTransaction(transaction)
	if token == nil {
		return
	}

	// set floor price of the collection
	token.FloorPrice, err = s.meClient.GetFloor(ctx, token.Symbol)
	if err != nil {
		log.Printf("meClient.GetFloor: %v", err)
		return
	}

	if token.Price < 1*token.FloorPrice {
		s.notificationService.SendNotification(ctx, &models.Action{Token: token})
	}

	if token.Type == "buy" || os.Getenv("ME_APIKEY") == "" || crypto.PrivateKey().String() == "" {
		return
	}

	// auto buy conditions
	if token.Price < 0.1 {
		signature, err := s.meClient.Buy(ctx, token)
		if err != nil {
			log.Println("Error while buying nft:", err.Error())
			return
		}
		log.Println("Signature -", signature)
		log.Println("Successfully bought item.")
	}
}

func (s *Service) parseTransaction(transaction *client.GetTransactionResponse) *models.Token {
	var (
		token *models.Token
		ok    bool
	)

	preTokenBalances := transaction.Meta.PreTokenBalances
	postTokenBalances := transaction.Meta.PostTokenBalances
	if len(preTokenBalances) == 0 {
		return nil
	}

	// Check if collections.json contains token
	mintAddress := preTokenBalances[0].Mint
	if token, ok = s.tokens[mintAddress]; !ok {
		return nil
	}

	price := getActionPrice(transaction.Meta.LogMessages)
	if price == 0 {
		return nil
	}

	actionType := getActionType(preTokenBalances[0].Owner, postTokenBalances[0].Owner)
	if actionType == "" {
		return nil
	}

	token.Type = actionType
	token.Timestamp = time.Now().UTC().Unix()
	token.BlockTimestamp = *transaction.BlockTime
	token.MintAddress = mintAddress
	token.Price = price
	token.TokenAddress = transaction.Transaction.Message.Accounts[2].String()
	token.Seller = transaction.Transaction.Message.Accounts[0].String()

	return token

}

func getActionPrice(logs []string) float64 {
	var price float64
	for _, msg := range logs {
		if strings.Contains(msg, "price") {
			if len(strings.Split(msg, "price\":")) < 2 {
				return 0
			}

			priceStr := strings.Split(strings.Split(msg, "price\":")[1], ",")[0]
			price, _ = strconv.ParseFloat(priceStr, 64)
			price /= 1e9
		}
	}

	return price
}

func getActionType(preTokenOwner, postTokenOwner string) string {
	if preTokenOwner == MEPublicKeyStr && postTokenOwner != MEPublicKeyStr {
		return "buy"
	} else if postTokenOwner == MEPublicKeyStr {
		return "list"
	}

	return ""
}
