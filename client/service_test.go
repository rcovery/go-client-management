package client_test

import (
	"context"
	"testing"

	"github.com/rcovery/go-client-management/client"
	"github.com/rcovery/go-client-management/client/mocks"
	"github.com/rcovery/go-client-management/client/postgres"
	infra_postgres "github.com/rcovery/go-client-management/internal/infra/postgres"
)

func newService() *client.Service {
	return client.NewService(&mocks.MockedRepository{}, &mocks.MockedCRM{})
}

func TestClientBusinessRules(t *testing.T) {
	t.Run("should create a valid client", func(t *testing.T) {
		ctx := context.Background()

		service := newService()

		clientData := &client.PostClientBody{
			Name:           "Ryan Test",
			Email:          "test@test.com",
			PortfolioValue: 100,
		}

		createdClient, creationErr := service.Insert(ctx, clientData)
		if creationErr != nil {
			t.Errorf("%v: %v", creationErr, clientData)
			t.FailNow()
		}

		if createdClient.Status != client.StatusPending {
			t.Errorf("status should be: %v. we received: %v", client.StatusPending, createdClient.Status)
			t.FailNow()
		}
	})

	t.Run("should not create a client with invalid email", func(t *testing.T) {
		ctx := context.Background()

		invalidEmail := "test@.com"

		service := newService()

		clientData := &client.PostClientBody{
			Name:           "Ryan Test",
			Email:          invalidEmail,
			PortfolioValue: 100,
		}

		_, creationErr := service.Insert(ctx, clientData)
		if creationErr == nil {
			t.Errorf("email \"%v\" should be invalid", invalidEmail)
			t.FailNow()
		}
	})

	t.Run("should not create a client with empty email", func(t *testing.T) {
		ctx := context.Background()

		service := newService()

		clientData := &client.PostClientBody{
			Name:           "Ryan Test",
			Email:          "",
			PortfolioValue: 100,
		}

		_, creationErr := service.Insert(ctx, clientData)
		if creationErr == nil {
			t.Errorf("we should not accept blank emails")
			t.FailNow()
		}
	})

	t.Run("should not create a client with empty name", func(t *testing.T) {
		ctx := context.Background()

		service := newService()

		clientData := &client.PostClientBody{
			Name:           "",
			Email:          "test@test.com",
			PortfolioValue: 100,
		}

		_, creationErr := service.Insert(ctx, clientData)
		if creationErr == nil {
			t.Errorf("we should not accept empty names")
			t.FailNow()
		}
	})

	t.Run("should not create a client with negative portfolio value", func(t *testing.T) {
		ctx := context.Background()

		service := newService()

		clientData := &client.PostClientBody{
			Name:           "Ryan Test",
			Email:          "test@test.com",
			PortfolioValue: -1,
		}

		_, creationErr := service.Insert(ctx, clientData)
		if creationErr == nil {
			t.Errorf("we should not accept negative portfolio values")
			t.FailNow()
		}
	})

	t.Run("should not create a client with zero portfolio value", func(t *testing.T) {
		ctx := context.Background()

		service := newService()

		clientData := &client.PostClientBody{
			Name:           "Ryan Test",
			Email:          "test@test.com",
			PortfolioValue: 0,
		}

		_, creationErr := service.Insert(ctx, clientData)
		if creationErr == nil {
			t.Errorf("we should not accept zero portfolio values")
			t.FailNow()
		}
	})

	t.Run("should update status and priority", func(t *testing.T) {
		ctx := context.Background()

		ClientID := client.ID("test-uuid")
		repo := &mocks.MockedRepository{
			StoredClient: &client.Client{
				ID:             &ClientID,
				Name:           "Update Test",
				Email:          "update@test.com",
				PortfolioValue: 250000,
			},
		}
		svc := client.NewService(repo, &mocks.MockedCRM{})

		updated, updateErr := svc.UpdateStatusAndPriority(ctx, "update@test.com", "card-123")
		if updateErr != nil {
			t.Fatalf("cannot update client: %v", updateErr)
		}

		if updated.Status != client.StatusProcessed {
			t.Errorf("status should be %v, got %v", client.StatusProcessed, updated.Status)
		}

		if updated.Priority == nil || *updated.Priority != client.HighPriority {
			t.Errorf("priority should be %v, got %v", client.HighPriority, updated.Priority)
		}

		if updated.PortfolioValue != 250000 {
			t.Errorf("portfolio value should remain %d, got %d", 250000, updated.PortfolioValue)
		}
	})

	t.Run("should set high priority for portfolio exactly 200000", func(t *testing.T) {
		ctx := context.Background()

		ClientID := client.ID("test-uuid-boundary")
		repo := &mocks.MockedRepository{
			StoredClient: &client.Client{
				ID:             &ClientID,
				Name:           "Boundary Test",
				Email:          "boundary@test.com",
				PortfolioValue: 200000,
			},
		}
		svc := client.NewService(repo, &mocks.MockedCRM{})

		updated, updateErr := svc.UpdateStatusAndPriority(ctx, "boundary@test.com", "card-123")
		if updateErr != nil {
			t.Fatalf("cannot update client: %v", updateErr)
		}

		if updated.Priority == nil || *updated.Priority != client.HighPriority {
			t.Errorf("priority should be %v, got %v", client.HighPriority, updated.Priority)
		}
	})

	t.Run("should set normal priority for portfolio 199999", func(t *testing.T) {
		ctx := context.Background()

		ClientID := client.ID("test-uuid-below")
		repo := &mocks.MockedRepository{
			StoredClient: &client.Client{
				ID:             &ClientID,
				Name:           "Below Test",
				Email:          "below@test.com",
				PortfolioValue: 199999,
			},
		}
		svc := client.NewService(repo, &mocks.MockedCRM{})

		updated, updateErr := svc.UpdateStatusAndPriority(ctx, "below@test.com", "card-123")
		if updateErr != nil {
			t.Fatalf("cannot update client: %v", updateErr)
		}

		if updated.Priority == nil || *updated.Priority != client.NormalPriority {
			t.Errorf("priority should be %v, got %v", client.NormalPriority, updated.Priority)
		}
	})
}

func TestClientCreation(t *testing.T) {
	t.Run("should create a valid client", func(t *testing.T) {
		ctx := context.Background()
		instance, postgresContainer := infra_postgres.SetupContainer(ctx, t)
		defer infra_postgres.TerminateContainer(postgresContainer)

		repo := postgres.NewRepository(instance)
		service := client.NewService(repo, &mocks.MockedCRM{})

		clientData := &client.PostClientBody{
			Name:           "Ryan Test",
			Email:          "test@test.com",
			PortfolioValue: 100,
		}

		createdClient, creationErr := service.Insert(ctx, clientData)
		if creationErr != nil {
			t.Errorf("cannot create client: %v", creationErr)
			t.FailNow()
		}

		if createdClient.Status != client.StatusPending {
			t.Errorf("status should be: %v. we received: %v", client.StatusPending, createdClient.Status)
			t.FailNow()
		}
	})

	t.Run("should select a client by email", func(t *testing.T) {
		ctx := context.Background()
		instance, postgresContainer := infra_postgres.SetupContainer(ctx, t)
		defer infra_postgres.TerminateContainer(postgresContainer)

		repo := postgres.NewRepository(instance)
		service := client.NewService(repo, &mocks.MockedCRM{})

		clientData := &client.PostClientBody{
			Name:           "Ryan Test",
			Email:          "test@test.com",
			PortfolioValue: 100,
		}

		createdClient, creationErr := service.Insert(ctx, clientData)
		if creationErr != nil {
			t.Errorf("cannot create client: %v", creationErr)
			t.FailNow()
		}

		foundClient, findErr := service.SelectByEmail(ctx, clientData.Email)
		if findErr != nil {
			t.Errorf("cannot find client by email: %v", findErr)
			t.FailNow()
		}
		if foundClient == nil {
			t.Errorf("client not found by email")
			t.FailNow()
		}
		if *foundClient.ID != *createdClient.ID {
			t.Errorf("found client ID doesn't match")
		}
	})
}
