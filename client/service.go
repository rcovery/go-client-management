package client

import (
	"context"
	"fmt"
	"log"
	"net/mail"
)

type Service struct {
	repo Repository
	crm  CRMGateway
}

func NewService(repo Repository, crm CRMGateway) *Service {
	return &Service{
		repo: repo,
		crm:  crm,
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
		log.Println(existingClientErr)
		return nil, fmt.Errorf("error.when.checking.existing.client")
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

	_, crmErr := s.crm.CreateCard(ctx, clientCreated, clientData.RequestType)
	if crmErr != nil {
		return nil, fmt.Errorf("could.not.update.card.in.crm")
	}

	return clientCreated, nil
}

func (s *Service) UpdateStatusAndPriority(ctx context.Context, email string, cardID string) (*Client, error) {
	parsedEmailAddress, emailErr := mail.ParseAddress(email)
	if emailErr != nil {
		return nil, fmt.Errorf("invalid.email")
	}

	existingClient, selectErr := s.repo.SelectByEmail(ctx, parsedEmailAddress.Address)
	if selectErr != nil {
		log.Println(selectErr)
		return nil, fmt.Errorf("error.when.checking.existing.client")
	}
	if existingClient == nil {
		return nil, fmt.Errorf("client.not.found")
	}

	priority := NormalPriority
	if existingClient.PortfolioValue >= 200000 {
		priority = HighPriority
	}

	existingClient.Status = StatusProcessed
	existingClient.Priority = &priority

	updateErr := s.repo.UpdateStatusAndPriority(ctx, existingClient)
	if updateErr != nil {
		return nil, fmt.Errorf("could.not.update.client")
	}

	_, crmErr := s.crm.UpdateCard(ctx, cardID, existingClient.Status, *existingClient.Priority)
	if crmErr != nil {
		return nil, fmt.Errorf("could.not.update.card.in.crm")
	}

	return existingClient, nil
}
