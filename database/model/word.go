package model

import (
	"context"
	"database/sql"
	"log"
)

type WordDto struct {
	ID            int64          `json:"id"`
	Word          string         `json:"word"`
	Transcription sql.NullString `json:"transcription"`
	Meaning       sql.NullString `json:"meaning"`
	Example       sql.NullString `json:"example"`
	WordLevel     sql.NullString `json:"word_level"`
	Translations  sql.NullString `json:"translations"`
	Frequency     sql.NullInt16  `json:"frequency"`
}

type WordStorage struct {
	db *sql.DB
}

func (s *WordStorage) GetWordById(ctx context.Context, wordId int64) (*WordDto, *DatabaseError) {
	// Implement getting logic here
	query := `SELECT id, word, transcription, meaning, example, word_level, translation FROM words WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryRowTimeout)
	defer cancel()

	wordDto := &WordDto{}
	err := s.db.QueryRowContext(
		ctx,
		query,
		wordId,
	).Scan(
		&wordDto.ID,
		&wordDto.Word,
		&wordDto.Transcription,
		&wordDto.Meaning,
		&wordDto.Example,
		&wordDto.WordLevel,
		&wordDto.Translations,
	)

	if err != nil {
		return nil, ProcessErrorFromDatabase(err, "GetWordById:45")
	}

	return wordDto, nil
}

func (s *WordStorage) GetWord(ctx context.Context, text string) (*WordDto, *DatabaseError) {
	// Implement getting logic here
	query := `SELECT id, word, transcription, meaning, example, word_level, translation FROM words WHERE word = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryRowTimeout)
	defer cancel()

	wordDto := &WordDto{}
	err := s.db.QueryRowContext(
		ctx,
		query,
		text,
	).Scan(
		&wordDto.ID,
		&wordDto.Word,
		&wordDto.Transcription,
		&wordDto.Meaning,
		&wordDto.Example,
		&wordDto.WordLevel,
		&wordDto.Translations,
	)

	if err != nil {
		return nil, ProcessErrorFromDatabase(err, "GetWord:73")
	}

	return wordDto, nil
}

func (s *WordStorage) SaveWords(ctx context.Context, wordDto *WordDto) (*WordDto, *DatabaseError) {
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
		return nil, ProcessErrorFromDatabase(err, "SaveWords:94")
	}

	return wordDto, nil
}

func (s *WordStorage) SaveWordWithBookConnection(ctx context.Context, book *BookDto, wordDto *WordDto) *DatabaseError {
	query := `INSERT INTO books_words (book_id, word_id, frequency) VALUES ($1, $2, $3)`
	log.Printf("book.ID : %d, wordDto.ID : %d", book.ID, wordDto.ID)

	_, err := s.db.ExecContext(
		ctx,
		query,
		book.ID,
		wordDto.ID,
		book.WordMap[wordDto.Word],
	)

	if err != nil {
		return ProcessErrorFromDatabase(err, "SaveWordWithBookConnection:112")
	}

	return nil
}

func (s *WordStorage) UpdateWords(ctx context.Context, wordDto *WordDto) (*WordDto, *DatabaseError) {
	// Implement updating logic here
	query := `UPDATE words SET word = $1, transcription = $2, meaning = $3, example = $4, word_level = $5, translations = $6 WHERE id = $7`

	_, err := s.db.ExecContext(
		ctx,
		query,
		wordDto.Word,
		wordDto.Transcription,
		wordDto.Meaning,
		wordDto.Example,
		wordDto.WordLevel,
		wordDto.Translations,
		wordDto.ID,
	)

	if err != nil {
		return nil, ProcessErrorFromDatabase(err, "UpdateWords:135")
	}

	return wordDto, nil
}
