package client

import "context"

type CRMWriter interface {
	CreateCard(ctx context.Context, clientData *Client, requestType string) (bool, error)
	UpdateCard(ctx context.Context, cardID string, status Status, priority Priority) (bool, error)
}

type CRMGateway interface {
	CRMWriter
}
