package webhook_test

import (
	"context"
	"testing"

	"github.com/rcovery/go-client-management/client"
	"github.com/rcovery/go-client-management/webhook"
	"github.com/rcovery/go-client-management/webhook/mocks"
)

func TestWebhookInsert(t *testing.T) {
	clientID := client.ID("test-client-uuid")

	t.Run("should insert valid webhook", func(t *testing.T) {
		repo := &mocks.MockedRepository{}
		clientReader := &mocks.MockedClientReader{
			StoredClient: &client.Client{ID: &clientID},
		}
		svc := webhook.NewService(repo, clientReader)

		body := &webhook.PostUpdatedCardBody{
			EventID: "evt-001", CardID: "card-123", ClienteEmail: "test@test.com",
		}
		w, err := svc.Insert(context.Background(), body)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if w == nil {
			t.Fatal("expected webhook, got nil")
		}
		if w.EventID != "evt-001" {
			t.Errorf("expected EventID 'evt-001', got '%s'", w.EventID)
		}
		if w.ClientID != clientID {
			t.Errorf("expected ClientID '%s', got '%s'", clientID, w.ClientID)
		}
	})

	t.Run("should block duplicate event_id", func(t *testing.T) {
		repo := &mocks.MockedRepository{
			StoredWebhook: &webhook.Webhook{EventID: "evt-001"},
		}
		clientReader := &mocks.MockedClientReader{
			StoredClient: &client.Client{ID: &clientID},
		}
		svc := webhook.NewService(repo, clientReader)

		body := &webhook.PostUpdatedCardBody{
			EventID: "evt-001", CardID: "card-123", ClienteEmail: "test@test.com",
		}
		_, err := svc.Insert(context.Background(), body)

		if err == nil {
			t.Fatal("expected error for duplicate event_id, got nil")
		}
		if err.Error() != "webhook.already.processed" {
			t.Errorf("expected 'webhook.already.processed', got '%s'", err.Error())
		}
	})

	t.Run("should reject empty event_id", func(t *testing.T) {
		svc := webhook.NewService(&mocks.MockedRepository{}, &mocks.MockedClientReader{})

		body := &webhook.PostUpdatedCardBody{
			EventID: "", CardID: "card-123", ClienteEmail: "test@test.com",
		}
		_, err := svc.Insert(context.Background(), body)

		if err == nil {
			t.Fatal("expected error for empty event_id, got nil")
		}
	})

	t.Run("should reject empty card_id", func(t *testing.T) {
		svc := webhook.NewService(&mocks.MockedRepository{}, &mocks.MockedClientReader{})

		body := &webhook.PostUpdatedCardBody{
			EventID: "evt-001", CardID: "", ClienteEmail: "test@test.com",
		}
		_, err := svc.Insert(context.Background(), body)

		if err == nil {
			t.Fatal("expected error for empty card_id, got nil")
		}
	})

	t.Run("should reject invalid email", func(t *testing.T) {
		svc := webhook.NewService(&mocks.MockedRepository{}, &mocks.MockedClientReader{})

		body := &webhook.PostUpdatedCardBody{
			EventID: "evt-001", CardID: "card-123", ClienteEmail: "invalido",
		}
		_, err := svc.Insert(context.Background(), body)

		if err == nil {
			t.Fatal("expected error for invalid email, got nil")
		}
	})

	t.Run("should reject when client not found", func(t *testing.T) {
		clientReader := &mocks.MockedClientReader{}
		svc := webhook.NewService(&mocks.MockedRepository{}, clientReader)

		body := &webhook.PostUpdatedCardBody{
			EventID: "evt-001", CardID: "card-123", ClienteEmail: "test@test.com",
		}
		_, err := svc.Insert(context.Background(), body)

		if err == nil {
			t.Fatal("expected error for client not found, got nil")
		}
		if err.Error() != "client.not.found" {
			t.Errorf("expected 'client.not.found', got '%s'", err.Error())
		}
	})
}
