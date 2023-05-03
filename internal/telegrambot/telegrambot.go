package telegrambot

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/fidesy/me-sniper/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	bot        *tgbotapi.BotAPI
	clientsIDs map[int64]bool
	mutex      sync.Mutex
}

func New(APIKey string) (*TelegramBot, error) {
	bot, err := tgbotapi.NewBotAPI(APIKey)
	if err != nil {
		return nil, err
	}

	return &TelegramBot{
		bot:        bot,
		clientsIDs: make(map[int64]bool),
	}, nil
}

func (tg *TelegramBot) Start(ctx context.Context, actions chan *models.Token) error {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := tg.bot.GetUpdatesChan(updateConfig)

	// listen new messages from clients
	go tg.handleUpdates(updates)

	// listen channel messages
	go tg.handleActions(actions)

	<-ctx.Done()
	return nil
}

func (tg *TelegramBot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			return
		}

		if update.Message.Text == "/start" {
			// add client to send action messages
			tg.mutex.Lock()
			defer tg.mutex.Unlock()
			tg.clientsIDs[update.Message.From.ID] = true
		}
	}
}

func (tg *TelegramBot) handleActions(actions chan *models.Token) {
	for action := range actions {
		action := action
		go func() {
			log.Println(action)
			// pretty string for output
			messageText := fmt.Sprintf("#%s \n%s \n<b>%s</b> \n%d/%d \n<b>%s for %.3fsol</b>\n<b>Floor: %.3fsol</b>  \n\nhttps://magiceden.io/item-details/%s", action.Symbol, action.Name, action.RarityStr, action.Rank, action.Supply, strings.ToUpper(action.Type), action.Price, action.FloorPrice, action.MintAddress)
			for clientID := range tg.clientsIDs {
				msg := tgbotapi.NewMessage(clientID, messageText)
				msg.ParseMode = "HTML"
				if _, err := tg.bot.Send(msg); err != nil {
					log.Println("error sending message:", err.Error())
				}
			}
		}()
	}
}
