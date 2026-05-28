package webhook

import (
	"time"

	"github.com/rcovery/go-client-management/client"
)

type Webhook struct {
	ID        *ID
	EventID   string
	CardID    string
	ClientID  client.ID
	CreatedAt time.Time
}

type PostUpdatedCardBody struct {
	EventID      string `json:"event_id"`
	CardID       string `json:"card_id"`
	ClienteEmail string `json:"cliente_email"`
	Timestamp    string `json:"timestamp"`
}
