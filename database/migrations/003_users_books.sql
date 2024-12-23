-- +goose Up
-- +goose StatementBegin
ALTER TABLE books ADD COLUMN user_id BIGSERIAL;
ALTER TABLE
    books
ADD
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE books DROP COLUMN user_id;
ALTER TABLE books DROP CONSTRAINT fk_user;
-- +goose StatementEnd
