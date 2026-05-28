package webhook

import (
	"context"
	"fmt"
	"net/mail"

	"github.com/rcovery/go-client-management/client"
)

type Service struct {
	repo         Repository
	clientReader client.Reader
}

func NewService(repo Repository, clientReader client.Reader) *Service {
	return &Service{
		repo:         repo,
		clientReader: clientReader,
	}
}

func (s *Service) Insert(ctx context.Context, body *PostUpdatedCardBody) (*Webhook, error) {
	if body.EventID == "" {
		return nil, fmt.Errorf("invalid.event_id")
	}

	if body.CardID == "" {
		return nil, fmt.Errorf("invalid.card_id")
	}

	parsedEmailAddress, emailErr := mail.ParseAddress(body.ClienteEmail)
	if emailErr != nil {
		return nil, fmt.Errorf("invalid.email")
	}

	existingWebhook, existingErr := s.repo.SelectByEventID(ctx, body.EventID)
	if existingErr != nil {
		return nil, fmt.Errorf("error.checking.existing.webhook")
	}
	if existingWebhook != nil {
		return nil, fmt.Errorf("webhook.already.processed")
	}

	existingClient, existingClientErr := s.clientReader.SelectByEmail(ctx, parsedEmailAddress.Address)
	if existingClientErr != nil {
		return nil, fmt.Errorf("error.checking.existing.client")
	}
	if existingClient == nil {
		return nil, fmt.Errorf("client.not.found")
	}

	id, idErr := NewID()
	if idErr != nil {
		return nil, fmt.Errorf("could.not.create.id")
	}

	webhookCreated := &Webhook{
		ID:       id,
		EventID:  body.EventID,
		CardID:   body.CardID,
		ClientID: *existingClient.ID,
	}

	_, insertErr := s.repo.Insert(ctx, webhookCreated)
	if insertErr != nil {
		return nil, fmt.Errorf("could.not.insert.webhook")
	}

	return webhookCreated, nil
}
