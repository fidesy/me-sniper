package crypto

import (
	"github.com/fidesy/me-sniper/internal/config"
	"github.com/gagliardetto/solana-go"
)

func PrivateKey() solana.PrivateKey {
	privateKeyString := config.Get(config.PrivateKey).(string)
	if privateKeyString == "" {
		return solana.PrivateKey{}
	}

	privateKey, err := solana.PrivateKeyFromBase58(privateKeyString)
	if err != nil {
		panic(err)
	}

	return privateKey
}
