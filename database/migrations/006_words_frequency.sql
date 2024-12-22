-- +goose Up
-- +goose StatementBegin
ALTER TABLE books_words ADD COLUMN frequency BIGSERIAL;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE books_words DROP COLUMN frequency;
-- +goose StatementEnd
