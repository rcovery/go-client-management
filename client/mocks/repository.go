package mocks

import (
	"context"

	"github.com/rcovery/go-client-management/client"
)

type MockedRepository struct{}

func (m *MockedRepository) Insert(ctx context.Context, clientData *client.Client) (*client.Client, error) {
	return clientData, nil
}

func (m *MockedRepository) SelectByEmail(ctx context.Context, email string) (*client.Client, error) {
	return nil, nil
}
