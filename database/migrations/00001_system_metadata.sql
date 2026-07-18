-- +goose Up
CREATE TABLE system_metadata (
  id UUID PRIMARY KEY,
  key TEXT NOT NULL UNIQUE,
  value JSONB NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

-- +goose Down
DROP TABLE system_metadata;
