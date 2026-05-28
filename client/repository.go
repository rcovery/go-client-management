package client

import "context"

type Writer interface {
	Insert(ctx context.Context, clientData *Client) (bool, error)
	UpdateStatusAndPriority(ctx context.Context, clientData *Client) error
}

type Reader interface {
	SelectByEmail(ctx context.Context, email string) (*Client, error)
}

type Repository interface {
	Reader
	Writer
}
