-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS data (
    id BIGINT,
    name VARCHAR(16) PRIMARY KEY,
    type VARCHAR(16),
    payload BYTEA
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE data;
-- +goose StatementEnd
