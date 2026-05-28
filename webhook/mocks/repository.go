package mocks

import (
	"context"

	"github.com/rcovery/go-client-management/client"
	"github.com/rcovery/go-client-management/webhook"
)

type MockedRepository struct {
	StoredWebhook *webhook.Webhook
}

func (m *MockedRepository) Insert(ctx context.Context, w *webhook.Webhook) (bool, error) {
	return true, nil
}

func (m *MockedRepository) SelectByEventID(ctx context.Context, eventID string) (*webhook.Webhook, error) {
	return m.StoredWebhook, nil
}

type MockedClientReader struct {
	StoredClient *client.Client
}

func (m *MockedClientReader) SelectByEmail(ctx context.Context, email string) (*client.Client, error) {
	if m.StoredClient == nil {
		return nil, nil
	}
	return m.StoredClient, nil
}
