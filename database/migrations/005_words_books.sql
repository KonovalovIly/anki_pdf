-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS books_words(
    book_id BIGSERIAL NOT NULL,
    word_id BIGSERIAL NOT NULL,
    CONSTRAINT book_id_word_id PRIMARY KEY (book_id, word_id)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE books_words;
-- +goose StatementEnd
