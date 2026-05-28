package mocks

import (
	"context"

	"github.com/rcovery/go-client-management/client"
)

type MockedRepository struct{}

func (m *MockedCRM) Insert(ctx context.Context, clientData *client.Client) (bool, error) {
	return true, nil
}

func (m *MockedCRM) SelectByEmail(ctx context.Context, email string) (*client.Client, error) {
	return nil, nil
}
