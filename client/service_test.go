package client

import (
	"context"
	"testing"
)

func TestNewClient(t *testing.T) {
	t.Run("should create a valid client", func(t *testing.T) {
		ctx := context.Background()

		clientData := &Client{
			Name:           "Ryan Test",
			Email:          "test@test.com",
			PortfolioValue: 100,
			RequestType:    "consulta",
		}

		createdClient, creationErr := Insert(ctx, clientData)
		if creationErr != nil {
			t.Errorf("%v: %v", creationErr, clientData)
			t.FailNow()
		}

		if createdClient.Status != StatusPending {
			t.Errorf("status should be: %v. we received: %v", StatusPending, createdClient.Status)
			t.FailNow()
		}
	})

	t.Run("should not create a client with invalid email", func(t *testing.T) {
		ctx := context.Background()

		invalidEmail := "test@.com"

		clientData := &Client{
			Name:           "Ryan Test",
			Email:          invalidEmail,
			PortfolioValue: 100,
			RequestType:    "consulta",
		}

		_, creationErr := Insert(ctx, clientData)
		if creationErr == nil {
			t.Errorf("email \"%v\" should be invalid", invalidEmail)
			t.FailNow()
		}
	})

	t.Run("should not create a client with empty email", func(t *testing.T) {
		ctx := context.Background()

		clientData := &Client{
			Name:           "Ryan Test",
			Email:          "",
			PortfolioValue: 100,
			RequestType:    "consulta",
		}

		_, creationErr := Insert(ctx, clientData)
		if creationErr == nil {
			t.Errorf("we should not accept blank emails")
			t.FailNow()
		}
	})

	t.Run("should not create a client with empty name", func(t *testing.T) {
		ctx := context.Background()

		clientData := &Client{
			Name:           "",
			Email:          "test@test.com",
			PortfolioValue: 100,
			RequestType:    "consulta",
		}

		_, creationErr := Insert(ctx, clientData)
		if creationErr == nil {
			t.Errorf("we should not accept empty names")
			t.FailNow()
		}
	})

	t.Run("should not create a client with negative portfolio value", func(t *testing.T) {
		ctx := context.Background()

		clientData := &Client{
			Name:           "Ryan Test",
			Email:          "test@test.com",
			PortfolioValue: -1,
			RequestType:    "consulta",
		}

		_, creationErr := Insert(ctx, clientData)
		if creationErr == nil {
			t.Errorf("we should not accept negative portfolio values")
			t.FailNow()
		}
	})

	t.Run("should not create a client with zero portfolio value", func(t *testing.T) {
		ctx := context.Background()

		clientData := &Client{
			Name:           "Ryan Test",
			Email:          "test@test.com",
			PortfolioValue: 0,
			RequestType:    "consulta",
		}

		_, creationErr := Insert(ctx, clientData)
		if creationErr == nil {
			t.Errorf("we should not accept zero portfolio value")
			t.FailNow()
		}
	})

	t.Run("should not create a client with empty request type", func(t *testing.T) {
		ctx := context.Background()

		clientData := &Client{
			Name:           "Ryan Test",
			Email:          "test@test.com",
			PortfolioValue: 100,
			RequestType:    "",
		}

		_, creationErr := Insert(ctx, clientData)
		if creationErr == nil {
			t.Errorf("we should not accept empty request type")
			t.FailNow()
		}
	})
}
