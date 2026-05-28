package handlers_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rcovery/go-client-management/client"
	clientMocks "github.com/rcovery/go-client-management/client/mocks"
	clientPostgres "github.com/rcovery/go-client-management/client/postgres"
	"github.com/rcovery/go-client-management/internal/http/handlers"
	infra_postgres "github.com/rcovery/go-client-management/internal/infra/postgres"
	"github.com/rcovery/go-client-management/webhook"
	webhookPostgres "github.com/rcovery/go-client-management/webhook/postgres"
)

func setupHandler(ctx context.Context, t *testing.T, seedClient *client.Client) (*handlers.Handler, *client.Service, *sql.DB, func()) {
	t.Helper()

	instance, postgresContainer := infra_postgres.SetupContainer(ctx, t)

	clientRepo := clientPostgres.NewRepository(instance)
	webhookRepo := webhookPostgres.NewRepository(instance)

	if seedClient != nil {
		_, err := clientRepo.Insert(ctx, seedClient)
		if err != nil {
			t.Fatalf("failed to seed client: %v", err)
		}
	}

	clientSvc := client.NewService(clientRepo, &clientMocks.MockedCRM{})
	webhookSvc := webhook.NewService(webhookRepo, clientSvc)
	h := handlers.NewHandler(clientSvc, webhookSvc)

	return h, clientSvc, instance, func() {
		infra_postgres.TerminateContainer(postgresContainer)
	}
}

func seedWebhook(ctx context.Context, t *testing.T, db *sql.DB, eventID string, clientID client.ID) {
	repo := webhookPostgres.NewRepository(db)
	id, err := webhook.NewID()
	if err != nil {
		t.Fatalf("failed to create webhook id: %v", err)
	}
	_, err = repo.Insert(ctx, &webhook.Webhook{
		ID:       id,
		EventID:  eventID,
		CardID:   "card-123",
		ClientID: clientID,
	})
	if err != nil {
		t.Fatalf("failed to seed webhook: %v", err)
	}
}

func TestHandleWebhookRequest(t *testing.T) {
	t.Run("should return high priority for portfolio", func(t *testing.T) {
		ctx := context.Background()
		clientID, err := client.NewID()
		if err != nil {
			t.Fatalf("failed to create client id: %v", err)
		}
		seedClient := &client.Client{
			ID:             clientID,
			Name:           "Test",
			Email:          "cliente@test.com",
			PortfolioValue: 200000,
		}

		h, clientSvc, _, cleanup := setupHandler(ctx, t, seedClient)
		defer cleanup()

		body, _ := json.Marshal(webhook.PostUpdatedCardBody{
			EventID: "evt-001", CardID: "card-123", ClienteEmail: "cliente@test.com",
		})
		req := httptest.NewRequest(http.MethodPost, "/webhooks/pipefy/card-updated", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		h.HandleWebhookRequest(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, rec.Code, rec.Body.String())
		}

		updatedClient, err := clientSvc.SelectByEmail(ctx, "cliente@test.com")
		if err != nil {
			t.Fatalf("failed to fetch updated client: %v", err)
		}
		if updatedClient.Status != client.StatusProcessed {
			t.Errorf("expected StatusProcessed, got %v", updatedClient.Status)
		}
		if updatedClient.Priority == nil || *updatedClient.Priority != client.HighPriority {
			t.Errorf("expected HighPriority, got %v", updatedClient.Priority)
		}
	})

	t.Run("should return normal priority for low portfolio", func(t *testing.T) {
		ctx := context.Background()
		clientID, err := client.NewID()
		if err != nil {
			t.Fatalf("failed to create client id: %v", err)
		}
		seedClient := &client.Client{
			ID:             clientID,
			Name:           "Test",
			Email:          "cliente@test.com",
			PortfolioValue: 100000,
		}

		h, clientSvc, _, cleanup := setupHandler(ctx, t, seedClient)
		defer cleanup()

		body, _ := json.Marshal(webhook.PostUpdatedCardBody{
			EventID: "evt-002", CardID: "card-123", ClienteEmail: "cliente@test.com",
		})
		req := httptest.NewRequest(http.MethodPost, "/webhooks/pipefy/card-updated", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		h.HandleWebhookRequest(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, rec.Code, rec.Body.String())
		}

		updatedClient, err := clientSvc.SelectByEmail(ctx, "cliente@test.com")
		if err != nil {
			t.Fatalf("failed to fetch updated client: %v", err)
		}
		if updatedClient.Status != client.StatusProcessed {
			t.Errorf("expected StatusProcessed, got %v", updatedClient.Status)
		}
		if updatedClient.Priority == nil || *updatedClient.Priority != client.NormalPriority {
			t.Errorf("expected NormalPriority, got %v", updatedClient.Priority)
		}
	})

	t.Run("should block duplicate event_id", func(t *testing.T) {
		ctx := context.Background()
		clientID, err := client.NewID()
		if err != nil {
			t.Fatalf("failed to create client id: %v", err)
		}
		seedClient := &client.Client{
			ID:             clientID,
			Name:           "Test",
			Email:          "cliente@test.com",
			PortfolioValue: 1000,
		}

		h, _, db, cleanup := setupHandler(ctx, t, seedClient)
		defer cleanup()

		seedWebhook(ctx, t, db, "evt-001", *clientID)

		body, _ := json.Marshal(webhook.PostUpdatedCardBody{
			EventID: "evt-001", CardID: "card-123", ClienteEmail: "cliente@test.com",
		})
		req := httptest.NewRequest(http.MethodPost, "/webhooks/pipefy/card-updated", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		h.HandleWebhookRequest(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
		}
		var resp map[string]string
		json.NewDecoder(rec.Body).Decode(&resp)
		if resp["message"] != "webhook.already.processed" {
			t.Errorf("expected message 'webhook.already.processed', got '%s'", resp["message"])
		}
	})

	t.Run("should reject when client not found", func(t *testing.T) {
		ctx := context.Background()

		h, _, _, cleanup := setupHandler(ctx, t, nil)
		defer cleanup()

		body, _ := json.Marshal(webhook.PostUpdatedCardBody{
			EventID: "evt-003", CardID: "card-123", ClienteEmail: "cliente@test.com",
		})
		req := httptest.NewRequest(http.MethodPost, "/webhooks/pipefy/card-updated", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		h.HandleWebhookRequest(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
		var resp map[string]string
		json.NewDecoder(rec.Body).Decode(&resp)
		if resp["error"] != "client.not.found" {
			t.Errorf("expected error 'client.not.found', got '%s'", resp["error"])
		}
	})

	t.Run("should reject empty body", func(t *testing.T) {
		ctx := context.Background()

		h, _, _, cleanup := setupHandler(ctx, t, nil)
		defer cleanup()

		req := httptest.NewRequest(http.MethodPost, "/webhooks/pipefy/card-updated", bytes.NewReader([]byte{}))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		h.HandleWebhookRequest(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("should reject invalid content type", func(t *testing.T) {
		ctx := context.Background()

		h, _, _, cleanup := setupHandler(ctx, t, nil)
		defer cleanup()

		body, _ := json.Marshal(webhook.PostUpdatedCardBody{
			EventID: "evt-005", CardID: "card-123", ClienteEmail: "cliente@test.com",
		})
		req := httptest.NewRequest(http.MethodPost, "/webhooks/pipefy/card-updated", bytes.NewReader(body))
		req.Header.Set("Content-Type", "text/plain")
		rec := httptest.NewRecorder()

		h.HandleWebhookRequest(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
		var resp map[string]string
		json.NewDecoder(rec.Body).Decode(&resp)
		if resp["error"] != "invalid.content.type" {
			t.Errorf("expected error 'invalid.content.type', got '%s'", resp["error"])
		}
	})

	t.Run("should reject method not allowed", func(t *testing.T) {
		ctx := context.Background()

		h, _, _, cleanup := setupHandler(ctx, t, nil)
		defer cleanup()

		req := httptest.NewRequest(http.MethodGet, "/webhooks/pipefy/card-updated", nil)
		rec := httptest.NewRecorder()

		h.HandleWebhookRequest(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
		}
	})
}
