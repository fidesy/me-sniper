package notificator

import (
	"context"
	"github.com/fidesy/me-sniper/internal/pkg/models"
	"log"
)

// Notificator just log all new listings
type Notificator struct{}

func New() *Notificator {
	return &Notificator{}
}

func (n *Notificator) SendNotification(_ context.Context, action *models.Action) {
	log.Println(action.String())
}
