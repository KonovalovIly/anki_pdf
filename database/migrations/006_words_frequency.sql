-- +goose Up
-- +goose StatementBegin
ALTER TABLE books_words ADD COLUMN frequency BIGSERIAL;
ALTER TABLE books_words ADD CONSTRAINT fk_book FOREIGN KEY (book_id) REFERENCES books(id);
ALTER TABLE books_words ADD CONSTRAINT fk_word FOREIGN KEY (word_id) REFERENCES words(id);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE books_words DROP COLUMN frequency;
ALTER TABLE books_words DROP CONSTRAINT fk_book;
ALTER TABLE books_words DROP CONSTRAINT fk_word;
-- +goose StatementEnd
