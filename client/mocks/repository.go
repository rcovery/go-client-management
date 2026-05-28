package mocks

import (
	"context"

	"github.com/rcovery/go-client-management/client"
)

type MockedRepository struct {
	StoredClient  *client.Client
	UpdatedClient *client.Client
}

func (m *MockedRepository) Insert(ctx context.Context, clientData *client.Client) (bool, error) {
	return true, nil
}

func (m *MockedRepository) SelectByEmail(ctx context.Context, email string) (*client.Client, error) {
	return m.StoredClient, nil
}

func (m *MockedRepository) UpdateStatusAndPriority(ctx context.Context, clientData *client.Client) error {
	m.UpdatedClient = clientData
	return nil
}
