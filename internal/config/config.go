package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	BotAPIToken = iota
	SolanaEndpoint
	MEAPIKey
	PrivateKey
)

var (
	ErrConfigNotFoundByKey = func(key int) error {
		return fmt.Errorf("config not found by key = %q", key)
	}
)

var conf *Config

type Config struct {
	BotAPIToken    string `yaml:"bot-api-token"`
	SolanaEndpoint string `yaml:"solana-endpoint"`
	MEAPIKey       string `yaml:"me-api-key"`
	PrivateKey     string `yaml:"private-key"`
}

func Init() error {
	body, err := os.ReadFile("./configs/values_local.yaml")
	if err != nil {
		return fmt.Errorf("os.ReadFile: %w", err)
	}

	err = yaml.Unmarshal(body, &conf)
	if err != nil {
		return fmt.Errorf("yaml.Unmarshal: %w", err)
	}

	return nil
}

func Get(key int) interface{} {
	switch key {
	case BotAPIToken:
		return conf.BotAPIToken
	case SolanaEndpoint:
		return conf.SolanaEndpoint
	case MEAPIKey:
		return conf.MEAPIKey
	case PrivateKey:
		return conf.PrivateKey
	default:
		panic(ErrConfigNotFoundByKey(key))
	}
}
