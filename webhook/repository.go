package webhook

import "context"

type Writer interface {
	Insert(ctx context.Context, w *Webhook) (bool, error)
}

type Reader interface {
	SelectByEventID(ctx context.Context, eventID string) (*Webhook, error)
}

type Repository interface {
	Reader
	Writer
}
