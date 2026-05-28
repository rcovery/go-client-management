package mocks

import (
	"context"

	"github.com/rcovery/go-client-management/client"
	"github.com/rcovery/go-client-management/webhook"
)

type MockedRepository struct{}

func (m *MockedRepository) Insert(ctx context.Context, w *webhook.Webhook) (bool, error) {
	return true, nil
}

func (m *MockedRepository) SelectByEventID(ctx context.Context, eventID string) (*webhook.Webhook, error) {
	return nil, nil
}

type MockedClientReader struct{}

func (m *MockedClientReader) SelectByEmail(ctx context.Context, email string) (*client.Client, error) {
	clientID := client.ID("test-client-uuid")
	return &client.Client{
		ID: &clientID,
	}, nil
}
