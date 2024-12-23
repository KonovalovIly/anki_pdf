
-- +goose Up
-- +goose StatementBegin
ALTER TABLE books DROP CONSTRAINT books_title_key;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE books DROP COLUMN books_title_key;
-- +goose StatementEnd
