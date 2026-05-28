package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/rcovery/go-client-management/client"
	"github.com/rcovery/go-client-management/client/pipefy"
	clientPostgres "github.com/rcovery/go-client-management/client/postgres"
	"github.com/rcovery/go-client-management/internal/config"
	"github.com/rcovery/go-client-management/webhook"
	webhookPostgres "github.com/rcovery/go-client-management/webhook/postgres"
)

type Handler struct {
	clientService  *client.Service
	webhookService *webhook.Service
}

func New(db *sql.DB) *Handler {
	pipefyGateway := pipefy.NewCRMGateway(config.GetString("PIPEFY_TOKEN"))

	clientRepo := clientPostgres.NewRepository(db)
	clientSvc := client.NewService(clientRepo, pipefyGateway)

	webhookRepo := webhookPostgres.NewRepository(db)
	webhookSvc := webhook.NewService(webhookRepo, clientSvc)

	return &Handler{
		clientService:  clientSvc,
		webhookService: webhookSvc,
	}
}

func NewHandler(clientSvc *client.Service, webhookSvc *webhook.Service) *Handler {
	return &Handler{
		clientService:  clientSvc,
		webhookService: webhookSvc,
	}
}

func writeJSONError(w http.ResponseWriter, status int, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	encodeErr := json.NewEncoder(w).Encode(map[string]string{
		"error": code,
	})
	if encodeErr != nil {
		log.Println("failed encoding json error:", encodeErr)
	}
}

const (
	KB int64 = 1024
	MB       = 1024 * KB
)
