package models

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type Token struct {
	Type           string  `json:"type"`
	Timestamp      int64   `json:"timestamp"`
	BlockTimestamp int64   `json:"blockTimestamp"`
	MintAddress    string  `json:"mintAddress"`
	Symbol         string  `json:"symbol"`
	Name           string  `json:"name"`
	Price          float64 `json:"price"`
	FloorPrice     float64 `json:"floorPrice"`
	RarityStr      string  `json:"rarityStr"`
	Rank           int     `json:"rank"`
	Supply         int     `json:"supply"`
	TokenAddress   string  `json:"tokenAddress"`
	Seller         string  `json:"seller"`
}

func (t *Token) String() string {
	return fmt.Sprintf(
		"Block timestamp:  %s \nNow: %s  \n#%s \n%s \n<b>%s</b> \n%d/%d \n<b>%s for %.3fsol</b>\n<b>Floor: %.3fsol</b>  \n\nhttps://magiceden.io/item-details/%s",
		time.Unix(t.BlockTimestamp, 0).Format("15:04:05 02-01-2006"),
		time.Unix(t.Timestamp, 0).Format("15:04:05 02-01-2006"),
		t.Symbol,
		t.Name,
		t.RarityStr,
		t.Rank,
		t.Supply,
		strings.ToUpper(t.Type),
		t.Price,
		t.FloorPrice,
		t.MintAddress,
	)
}

func LoadTokens() (map[string]*Token, error) {
	var collections map[string]*Token
	data, err := os.ReadFile("./data/collections.json")
	if err != nil {
		return nil, fmt.Errorf("os.ReadFlie: %w", err)
	}

	err = json.Unmarshal(data, &collections)
	if err != nil {
		return nil, err
	}

	return collections, nil
}
