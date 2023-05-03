package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fidesy/me-sniper/internal/models"
	"github.com/fidesy/me-sniper/internal/sniper"
	"github.com/fidesy/me-sniper/internal/telegrambot"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	checkError(err)

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
		syscall.SIGQUIT,
	)
	defer cancel()

	var actions = make(chan *models.Token)

	// create sniper instance
	s, err := sniper.New(&sniper.Options{
		Endpoint:   os.Getenv("NODE_ENDPOINT"),
		Actions:    actions,
		PrivateKey: os.Getenv("PRIVATE_KEY"),
	})
	checkError(err)

	// run sniper concurrently
	go func() {
		err = s.Start(ctx)
		checkError(err)
	}()

	telegramAPIKey := os.Getenv("TELEGRAM_APIKEY")
	if telegramAPIKey != "" {
		// create and start telegram bot
		tgBot, err := telegrambot.New(telegramAPIKey)
		checkError(err)

		err = tgBot.Start(ctx, actions)
		checkError(err)
	} else {
		// just logs
		for action := range actions {
			action := action

			go func() {
				log.Println(action)
			}()
		}
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
