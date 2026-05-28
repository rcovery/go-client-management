package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/rcovery/go-client-management/webhook"
)

func (h *Handler) HandleWebhook() {
	http.HandleFunc("/webhooks/pipefy/card-updated", h.HandleWebhookRequest)
}

func (h *Handler) HandleWebhookRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx, ctxCancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer ctxCancel()

	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		writeJSONError(w, http.StatusBadRequest, "invalid.content.type")
		return
	}

	rawBody := http.MaxBytesReader(w, r.Body, 2*MB)
	body, err := io.ReadAll(rawBody)
	if err != nil {
		log.Println("failed reading body:", err)
		writeJSONError(w, http.StatusBadRequest, "invalid.body")
		return
	}
	if len(body) == 0 {
		writeJSONError(w, http.StatusBadRequest, "empty.body")
		return
	}

	var webhookBody webhook.PostUpdatedCardBody
	err = json.Unmarshal(body, &webhookBody)
	if err != nil {
		log.Println("failed decoding json:", err)
		writeJSONError(w, http.StatusBadRequest, "invalid.json")
		return
	}

	_, creationErr := h.webhookService.Insert(ctx, &webhookBody)
	if creationErr != nil {
		log.Println(creationErr)
		writeJSONError(w, http.StatusBadRequest, creationErr.Error())
		return
	}

	_, clientUpdateErr := h.clientService.UpdateStatusAndPriority(ctx, webhookBody.ClienteEmail, webhookBody.CardID)
	if clientUpdateErr != nil {
		log.Println(clientUpdateErr)
		writeJSONError(w, http.StatusBadRequest, clientUpdateErr.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	enc := json.NewEncoder(w)
	enc.Encode(map[string]any{
		"success": true,
		"message": "created.successfully",
	})
}
