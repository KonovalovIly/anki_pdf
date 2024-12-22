-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS books_words(
    book_id BIGSERIAL,
    word_id BIGSERIAL,
    CONSTRAINT book_id_word_id PRIMARY KEY (book_id, word_id),
    CONSTRAINT fk_book_id FOREIGN KEY (book_id) REFERENCES books(id),
    CONSTRAINT fk_word_id FOREIGN KEY (word_id) REFERENCES words(id)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE books_words;
-- +goose StatementEnd
