package client

import "context"

type Writer interface {
	Insert(ctx context.Context, clientData *Client) (*Client, error)
}

type Reader interface {
	SelectByEmail(ctx context.Context, email string) (*Client, error)
	SelectByIdempotencyKey(ctx context.Context, idempotencyKey IdempotencyKey) (*Client, error)
}

type Repository struct {
	Reader
	Writer
}
