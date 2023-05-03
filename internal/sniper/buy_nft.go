package sniper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gagliardetto/solana-go"

	"github.com/fidesy/me-sniper/internal/models"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go/rpc"
)

var (
	buyURL = "https://api-mainnet.magiceden.dev/v2/instructions/buy_now?buyer=%s&seller=%s&auctionHouseAddress=E8cU1WiRWjanGxmn96ewBgk9vPTcL6AEZ1t6F6fkgUWe&tokenMint=%s&tokenATA=%s&price=%f&sellerExpiry=-1&useV2=false&buyerCreatorRoyaltyPercent=0"
)

func (s *Sniper) BuyNFT(token *models.Token) (string, error) {
	url := fmt.Sprintf(buyURL,
		s.privateKey.PublicKey(), token.Seller, token.MintAddress, token.TokenAddress, token.Price)

	txSigned, err := getTransactionData(url)
	if err != nil {
		return "", err
	}

	transaction, err := solana.TransactionFromDecoder(bin.NewBorshDecoder(txSigned))
	if err != nil {
		return "", err
	}

	cli := rpc.New(os.Getenv("NODE_ENDPOINT"))

	// Sign transaction message
	messageContent, _ := transaction.Message.MarshalBinary()
	sign, _ := s.privateKey.Sign(messageContent)
	transaction.Signatures[0] = sign

	signature, err := cli.SendTransactionWithOpts(context.TODO(), transaction, rpc.TransactionOpts{SkipPreflight: true})
	if err != nil {
		return "", err
	}

	return signature.String(), nil
}

func getTransactionData(url string) ([]byte, error) {
	client := http.Client{}

	req, _ := http.NewRequest("GET", url, nil)
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

	var response models.BuyResponse
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response.TxSigned.Data, nil
}
