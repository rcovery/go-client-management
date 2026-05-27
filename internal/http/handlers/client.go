package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/rcovery/go-client-management/client"
	"github.com/rcovery/go-client-management/client/postgres"
)

func HandleClient(DB *sql.DB) {
	repo := postgres.NewRepository(DB)
	service := client.NewService(repo)

	http.HandleFunc("/clientes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			{
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

				var clientBody client.PostClientBody
				err = json.Unmarshal(body, &clientBody)
				if err != nil {
					log.Println("failed decoding json:", err)
					writeJSONError(w, http.StatusBadRequest, "invalid.json")
					return
				}

				createdClient, creationErr := service.Insert(ctx, &clientBody)
				if creationErr != nil || createdClient == nil {
					log.Println(creationErr)
					writeJSONError(w, http.StatusBadRequest, string(creationErr.Error()))
					break
				}

				enc := json.NewEncoder(w)
				enc.Encode(map[string]any{
					"success": true,
					"message": "created.successfully",
				})
				break
			}
		default:
			{
				writeJSONError(w, http.StatusMethodNotAllowed, "method.not.allowed")
			}
		}
	})
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
