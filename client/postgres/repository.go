package postgres

import (
	"context"
	"database/sql"

	"github.com/rcovery/go-client-management/client"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(DB *sql.DB) *Repository {
	return &Repository{
		DB: DB,
	}
}

func (r *Repository) SelectByEmail(ctx context.Context, email string) (*client.Client, error) {
	query := `SELECT id, name, email, portfolio_value, status FROM clients WHERE email = $1`

	var idStr string
	clientFound := &client.Client{}

	err := r.DB.QueryRowContext(ctx, query, email).Scan(&idStr, &clientFound.Name, &clientFound.Email, &clientFound.PortfolioValue, &clientFound.Status)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	clientFound.ID = (*client.ID)(&idStr)

	return clientFound, nil
}

func (r *Repository) Insert(ctx context.Context, clientData *client.Client) (bool, error) {
	query := `INSERT INTO clients (id, name, email, portfolio_value, status) VALUES ($1, $2, $3, $4, $5)`

	_, execErr := r.DB.ExecContext(ctx, query,
		string(*clientData.ID),
		clientData.Name,
		clientData.Email,
		clientData.PortfolioValue,
		string(clientData.Status),
	)
	if execErr != nil {
		return false, execErr
	}

	return true, nil
}
