package app

import (
	"context"
	"fmt"
	"github.com/fidesy/me-sniper/internal/config"
	"github.com/fidesy/me-sniper/internal/pkg/magiceden"
	"github.com/fidesy/me-sniper/internal/pkg/sniper"
	"github.com/fidesy/me-sniper/internal/pkg/telegrambot"
	"golang.org/x/sync/errgroup"
)

func Run(ctx context.Context) error {
	meClient := magiceden.New(ctx)

	telegramBot, err := telegrambot.New(
		config.Get(config.BotAPIToken).(string),
	)
	if err != nil {
		return fmt.Errorf("telegrambot.New: %w", err)
	}

	sniperService, err := sniper.New(
		meClient,
		telegramBot,
	)
	if err != nil {
		return fmt.Errorf("sniper.New: %w", err)
	}

	errGroup, ctx := errgroup.WithContext(ctx)

	errGroup.Go(func() error {
		if err = telegramBot.Run(ctx); err != nil {
			return fmt.Errorf("telegramBot.Run: %w", err)
		}

		return nil
	})

	errGroup.Go(func() error {
		if err = sniperService.Run(ctx); err != nil {
			return fmt.Errorf("sniper.Run: %w", err)
		}

		return nil
	})

	return errGroup.Wait()
}
