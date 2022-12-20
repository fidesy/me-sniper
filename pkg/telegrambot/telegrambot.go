package telegrambot

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/fidesy/me-sniper/pkg/models"
)

type TelegramBot struct {
	bot     *tgbotapi.BotAPI
	actions chan *models.Token
}

func New(API_KEY string, actions chan *models.Token) (*TelegramBot, error) {
	bot, err := tgbotapi.NewBotAPI(API_KEY)
	if err != nil {
		return nil, err
	}

	return &TelegramBot{
		bot:     bot,
		actions: actions,
	}, nil
}

func (tg *TelegramBot) Start() error {
	var clientsIDs = make(map[int64]bool)
	// listen new messages from clients
	go func() {
		updateConfig := tgbotapi.NewUpdate(0)
		updateConfig.Timeout = 30
		updates := tg.bot.GetUpdatesChan(updateConfig)
		for update := range updates {
			if update.Message == nil {
				return
			}

			if update.Message.Text == "/start" {
				// add client to send action messages
				clientsIDs[update.Message.From.ID] = true
			}
		}
	}()

	// listen channel messages
	for action := range tg.actions {
		action := action
		go func() {
			log.Println(action)
			// pretty string for output
			messageText := fmt.Sprintf("#%s \n%s \n<b>%s</b> \n%d/%d \n<b>%s for %.3fsol</b>\n<b>Floor: %.3fsol</b>  \n\nhttps://magiceden.io/item-details/%s", action.Symbol, action.Name, action.RarityStr, action.Rank, action.Supply, strings.ToUpper(action.Type), action.Price, action.FloorPrice, action.MintAddress)
			for clientID := range clientsIDs {
				msg := tgbotapi.NewMessage(clientID, messageText)
				msg.ParseMode = "HTML"
				_, err := tg.bot.Send(msg)
				if err != nil {
					log.Fatal(err)
				}
			}
		}()
	}

	return nil
}
