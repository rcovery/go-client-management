package mocks

import (
	"context"

	"github.com/rcovery/go-client-management/client"
)

type MockedCRM struct{}

func (m *MockedCRM) CreateCard(ctx context.Context, clientData *client.Client, requestType string) (bool, error) {
	return true, nil
}

func (m *MockedCRM) UpdateCard(ctx context.Context, cardID string, status client.Status, priority client.Priority) (bool, error) {
	return true, nil
}
