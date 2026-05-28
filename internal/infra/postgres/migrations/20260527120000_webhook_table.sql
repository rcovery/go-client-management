-- +goose Up
CREATE TABLE webhooks (
  id UUID PRIMARY KEY,
  event_id TEXT NOT NULL,
  card_id TEXT NOT NULL,
  client_id UUID NOT NULL REFERENCES clients(id),
  created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_webhooks_event_id ON webhooks (event_id);

-- +goose Down
DROP TABLE IF EXISTS webhooks;
