-- +goose Up
-- +goose StatementBegin
ALTER TABLE books ADD COLUMN words_count INTEGER;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE books DROP COLUMN words_count;
-- +goose StatementEnd
