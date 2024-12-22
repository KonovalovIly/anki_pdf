-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS words(
    id BIGSERIAL PRIMARY KEY,
    word VARCHAR(45) UNIQUE NOT NULL,
    transcription VARCHAR(100),
    meaning VARCHAR(255),
    example VARCHAR(255),
    word_level VARCHAR(2),
    translation VARCHAR(100)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE words;
-- +goose StatementEnd
