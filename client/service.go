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

func (s *Service) SelectByEmail(ctx context.Context, email string) (*Client, error) {
	parsedEmailAddress, emailErr := mail.ParseAddress(email)
	if emailErr != nil {
		return nil, fmt.Errorf("invalid.email")
	}

	return s.repo.SelectByEmail(ctx, parsedEmailAddress.Address)
}

func (s *Service) Insert(ctx context.Context, clientData *PostClientBody) (*Client, error) {
	if clientData.Name == "" {
		return nil, fmt.Errorf("invalid.name")
	}

	if clientData.PortfolioValue <= 0 {
		return nil, fmt.Errorf("invalid.portfolio.value")
	}

	parsedEmailAddress, emailErr := mail.ParseAddress(clientData.Email)
	if emailErr != nil {
		return nil, fmt.Errorf("invalid.email")
	}

	alreadyExistingClient, existingClientErr := s.repo.SelectByEmail(ctx, parsedEmailAddress.Address)
	if existingClientErr != nil {
		return nil, fmt.Errorf("error.checking.existing.client")
	}
	if alreadyExistingClient != nil {
		return nil, fmt.Errorf("client.already.exists")
	}

	id, idErr := NewID()
	if idErr != nil {
		return nil, fmt.Errorf("could.not.create.id")
	}

	clientCreated := &Client{
		Name:           clientData.Name,
		PortfolioValue: clientData.PortfolioValue,
		ID:             id,
		Email:          parsedEmailAddress.Address,
		Status:         StatusPending,
	}

	_, insertErr := s.repo.Insert(ctx, clientCreated)
	if insertErr != nil {
		return nil, fmt.Errorf("could.not.insert.client")
	}

	return clientCreated, nil
}
