package model

import (
	"context"
	"database/sql"
)

type WordDto struct {
	ID            int    `json:"id"`
	Word          string `json:"word"`
	Transcription string `json:"transcription"`
	Meaning       string `json:"meaning"`
	Example       string `json:"example"`
	WordLevel     string `json:"word_level"`
	Translations  string `json:"translations"`
}

type WordStorage struct {
	db *sql.DB
}

func (s *WordStorage) SaveWords(ctx context.Context, book *BookDto, wordDto *WordDto) (*WordDto, error) {
	// Implement saving logic here
	query := `INSERT INTO words (word) VALUES ($1) RETURNING id`
	ctx, cancel := context.WithTimeout(ctx, QueryRowTimeout)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		wordDto.Word,
	).Scan(
		&wordDto.ID,
	)

	if err != nil {
		return nil, err
	}

	query = `INSERT INTO books_words (book_id, word_id, frequency) VALUES ($1, $2, $3)`
	_, err = s.db.ExecContext(
		ctx,
		query,
		book.ID,
		wordDto.ID,
		book.WordMap[wordDto.Word],
	)
	if err != nil {
		return nil, err
	}

	return wordDto, nil
}
