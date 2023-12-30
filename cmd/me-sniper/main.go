package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/fidesy/me-sniper/internal/app"
	"github.com/fidesy/me-sniper/internal/config"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGTERM,
	)
	defer cancel()

	err := config.Init()
	if err != nil {
		log.Fatalf("config.Init: %v", err)
	}

	if err = app.Run(ctx); err != nil {
		log.Fatalf("service.Run: %v", err)
	}
}
