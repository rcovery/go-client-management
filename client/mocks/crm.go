package mocks

import (
	"context"

	"github.com/rcovery/go-client-management/client"
)

type MockedCRM struct{}

func (m *MockedCRM) CreateCard(ctx context.Context, clientData *client.Client) (bool, error) {
	return true, nil
}
