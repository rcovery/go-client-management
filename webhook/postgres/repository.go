package postgres

import (
	"context"
	"database/sql"

	"github.com/rcovery/go-client-management/client"
	"github.com/rcovery/go-client-management/webhook"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(DB *sql.DB) *Repository {
	return &Repository{
		DB: DB,
	}
}

func (r *Repository) Insert(ctx context.Context, w *webhook.Webhook) (bool, error) {
	query := `INSERT INTO webhooks (id, event_id, card_id, client_id) VALUES ($1, $2, $3, $4)`

	_, execErr := r.DB.ExecContext(ctx, query,
		string(*w.ID),
		w.EventID,
		w.CardID,
		string(w.ClientID),
	)
	if execErr != nil {
		return false, execErr
	}

	return true, nil
}

func (r *Repository) SelectByEventID(ctx context.Context, eventID string) (*webhook.Webhook, error) {
	query := `SELECT id, event_id, card_id, client_id, created_at FROM webhooks WHERE event_id = $1`

	var idStr string
	var clientIDStr string
	w := &webhook.Webhook{}

	err := r.DB.QueryRowContext(ctx, query, eventID).Scan(&idStr, &w.EventID, &w.CardID, &clientIDStr, &w.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	w.ID = (*webhook.ID)(&idStr)
	w.ClientID = client.ID(clientIDStr)

	return w, nil
}
