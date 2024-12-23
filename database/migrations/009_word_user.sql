-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users_words(
    user_id BIGSERIAL NOT NULL,
    word_id BIGSERIAL NOT NULL,
    is_learned BOOLEAN NOT NULL,
    CONSTRAINT user_id_word_id PRIMARY KEY (user_id, word_id),
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_word FOREIGN KEY (word_id) REFERENCES words(id)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users_words;
-- +goose StatementEnd
