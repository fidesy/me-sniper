package magiceden

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fidesy/me-sniper/internal/config"
	"github.com/fidesy/me-sniper/internal/pkg/cache"
	"github.com/fidesy/me-sniper/internal/pkg/crypto"
	"github.com/fidesy/me-sniper/internal/pkg/models"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	TTL                = 30 * time.Second
	FloorPriceEndpoint = "https://api-mainnet.magiceden.dev/v2/collections/%s/stats"
	BuyEndpoint        = "https://api-mainnet.magiceden.dev/v2/instructions/buy_now?buyer=%s&seller=%s&auctionHouseAddress=E8cU1WiRWjanGxmn96ewBgk9vPTcL6AEZ1t6F6fkgUWe&tokenMint=%s&tokenATA=%s&price=%f&sellerExpiry=-1&useV2=false&buyerCreatorRoyaltyPercent=0"
)

type Service struct {
	floorPriceCache *cache.Cache
	client          *http.Client
}

func New(ctx context.Context) *Service {
	return &Service{
		floorPriceCache: cache.New(ctx, TTL),
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (s *Service) GetFloor(ctx context.Context, symbol string) (float64, error) {
	floorPrice := float64(0)
	exists := s.floorPriceCache.Get(symbol, &floorPrice)
	if exists {
		return floorPrice, nil
	}

	request, _ := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(FloorPriceEndpoint, symbol), nil)
	resp, err := s.client.Do(request)
	if err != nil {
		return 0, fmt.Errorf("client.Do: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("response status is not success")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("io.ReadAll: %w", err)
	}

	var floorResp FloorResponse
	if err = json.Unmarshal(body, &floorResp); err != nil {
		return 0, fmt.Errorf("json.Unmarshal: %w", err)
	}

	floorPrice = floorResp.FloorPrice / 1e9

	s.floorPriceCache.Set(symbol, floorPrice, TTL)

	return floorPrice, nil
}

func (s *Service) Buy(ctx context.Context, token *models.Token) (string, error) {
	privateKey := crypto.PrivateKey()

	url := fmt.Sprintf(
		BuyEndpoint,
		privateKey.PublicKey(),
		token.Seller,
		token.MintAddress,
		token.TokenAddress,
		token.Price,
	)

	txSigned, err := s.transactionData(ctx, url)
	if err != nil {
		return "", err
	}

	transaction, err := solana.TransactionFromDecoder(bin.NewBorshDecoder(txSigned))
	if err != nil {
		return "", err
	}

	cli := rpc.New(config.Get(config.SolanaEndpoint).(string))

	// Sign transaction message
	messageContent, _ := transaction.Message.MarshalBinary()
	sign, _ := privateKey.Sign(messageContent)
	transaction.Signatures[0] = sign

	signature, err := cli.SendTransactionWithOpts(ctx, transaction, rpc.TransactionOpts{SkipPreflight: true})
	if err != nil {
		return "", err
	}

	return signature.String(), nil
}

func (s *Service) transactionData(ctx context.Context, url string) ([]byte, error) {
	client := http.Client{}

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	req.Header.Add("Authorization", "Bearer "+os.Getenv("ME_APIKEY"))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, errors.New(string(body))
	}

	var response BuyResponse
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response.TxSigned.Data, nil
}
