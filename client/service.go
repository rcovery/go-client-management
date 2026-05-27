package client

import (
	"context"
	"fmt"
	"net/mail"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func Insert(ctx context.Context, clientData *Client) (*Client, error) {
	// if clientData.ID == nil {
	// }

	parsedEmailAddress, emailErr := mail.ParseAddress(clientData.Email)
	if emailErr != nil {
		return nil, fmt.Errorf("invalid.email")
	}

	id, idErr := NewID()
	if idErr != nil {
		return nil, fmt.Errorf("could.not.create.id")
	}

	clientData.ID = id
	clientData.Email = parsedEmailAddress.Address
	clientData.Status = StatusPending

	return clientData, nil
}
