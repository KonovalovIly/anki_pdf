-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS books(
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) UNIQUE NOT NULL,
    added_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE books;
-- +goose StatementEnd
