package client

import "context"

type CRMWriter interface {
	CreateCard(ctx context.Context, clientData *Client) (bool, error)
}

type CRMGateway interface {
	CRMWriter
}
