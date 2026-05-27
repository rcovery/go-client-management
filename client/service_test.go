package client

import (
	"context"
	"testing"
)

func TestNewClient(t *testing.T) {
	t.Run("should create a valid client", func(t *testing.T) {
		ctx := context.Background()

		clientData := &Client{
			Email: "test@test.com",
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
			Email: invalidEmail,
		}

		_, creationErr := Insert(ctx, clientData)
		if creationErr == nil {
			t.Errorf("email \"%v\" should be invalid", invalidEmail)
			t.FailNow()
		}
	})

	t.Run("should not create a client with blank email", func(t *testing.T) {
		ctx := context.Background()

		invalidEmail := ""

		clientData := &Client{
			Email: invalidEmail,
		}

		_, creationErr := Insert(ctx, clientData)
		if creationErr == nil {
			t.Errorf("we should not accept blank emails")
			t.FailNow()
		}
	})
}
