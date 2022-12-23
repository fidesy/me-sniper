package sniper

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/fidesy/me-sniper/pkg/models"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/portto/solana-go-sdk/client"
)

func GetTransactionData(url string) ([]byte, error) {
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

func BuyNFT(c *client.Client, privateKey solana.PrivateKey, url string) (string, error) {
	txSigned, err := GetTransactionData(url)
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
	s, _ := privateKey.Sign(messageContent)
	transaction.Signatures[0] = s

	signature, err := cli.SendTransactionWithOpts(context.TODO(), transaction, rpc.TransactionOpts{SkipPreflight: true})
	if err != nil {
		return "", err
	}

	return signature.String(), nil
}
