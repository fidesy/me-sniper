package main

import (
	"log"
	"os"

	"github.com/fidesy/me-sniper/internal/models"
	"github.com/fidesy/me-sniper/internal/sniper"
	"github.com/fidesy/me-sniper/internal/telegrambot"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	checkError(err)

	var actions = make(chan *models.Token, 5)

	// create sniper instance
	s, err := sniper.New(&sniper.Options{
		Endpoint:   os.Getenv("NODE_ENDPOINT"),
		Actions:    actions,
		PrivateKey: os.Getenv("PRIVATE_KEY"),
	})
	checkError(err)

	// run sniper concurrently
	go func() {
		err = s.Start()
		checkError(err)
	}()

	telegramAPIKey := os.Getenv("TELEGRAM_APIKEY")
	if telegramAPIKey != "" {
		// create and start telegram bot
		tgBot, err := telegrambot.New(telegramAPIKey, actions)
		checkError(err)

		err = tgBot.Start()
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
