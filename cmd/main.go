package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/fidesy/me-sniper/pkg/models"
	"github.com/fidesy/me-sniper/pkg/sniper"
	"github.com/fidesy/me-sniper/pkg/telegrambot"
)

func main() {
	err := godotenv.Load()
	checkError(err)

	var actions = make(chan *models.Token, 5)

	// create sniper instance
	s, err := sniper.New(os.Getenv("NODE_ENDPOINT"), actions)
	checkError(err)

	go func() {
		err = s.Start()
		checkError(err)
	}()

	TELEGRAM_APIKEY := os.Getenv("TELEGRAM_APIKEY")
	if TELEGRAM_APIKEY != "" {
		// create and start telegram bot
		tgbot, err := telegrambot.New(TELEGRAM_APIKEY, actions)
		checkError(err)

		err = tgbot.Start()
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
