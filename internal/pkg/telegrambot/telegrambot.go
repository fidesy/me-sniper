package telegrambot

import (
	"context"
	"log"
	"sync"

	"github.com/fidesy/me-sniper/internal/pkg/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type (
	Service struct {
		bot        *tgbotapi.BotAPI
		clientsIDs map[int64]bool
		mutex      sync.Mutex
	}
)

func New(apiKey string) (*Service, error) {
	bot, err := tgbotapi.NewBotAPI(apiKey)
	if err != nil {
		return nil, err
	}

	return &Service{
		bot:        bot,
		clientsIDs: make(map[int64]bool),
	}, nil
}

func (s *Service) Run(ctx context.Context) error {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := s.bot.GetUpdatesChan(updateConfig)

	for {
		select {
		case <-ctx.Done():
			return nil
		case update := <-updates:
			go s.handleUpdate(update)
		}
	}
}

func (s *Service) handleUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	if update.Message.Text == "/start" {
		// add client to send action messages
		s.mutex.Lock()
		s.clientsIDs[update.Message.From.ID] = true
		s.mutex.Unlock()
	}
}

func (s *Service) SendNotification(_ context.Context, action *models.Action) {
	log.Println(action)
	for clientID := range s.clientsIDs {
		msg := tgbotapi.NewMessage(clientID, action.String())
		msg.ParseMode = "HTML"
		if _, err := s.bot.Send(msg); err != nil {
			log.Println("error sending message:", err.Error())
		}
	}
}
