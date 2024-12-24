package storage

import (
	"context"
	"database/sql"

	database_models "github.com/KonovalovIly/anki_pdf/database/model"
)

type WordStorage struct {
	db *sql.DB
}

func (s *WordStorage) GetWordById(ctx context.Context, wordId int64) (*database_models.WordDto, *database_models.DatabaseError) {
	// Implement getting logic here
	query := `SELECT
		id,
		word,
		transcription,
		meaning,
		example,
		word_level,
		translation
		FROM words WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryRowTimeout)
	defer cancel()

	wordDto := &database_models.WordDto{}
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
		&wordDto.Translation,
	)

	if err != nil {
		return nil, database_models.ProcessErrorFromDatabase(err, "GetWordById:QueryRowContext")
	}

	return wordDto, nil
}

func (s *WordStorage) GetWordByName(ctx context.Context, text string) (*database_models.WordDto, *database_models.DatabaseError) {
	query := `SELECT
		id,
		word,
		transcription,
		meaning,
		example,
		word_level,
		translation
		FROM words WHERE word = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryRowTimeout)
	defer cancel()

	wordDto := &database_models.WordDto{}
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
		&wordDto.Translation,
	)

	if err != nil {
		return nil, database_models.ProcessErrorFromDatabase(err, "GetWord:QueryRowContext")
	}

	return wordDto, nil
}

func (s *WordStorage) SaveWord(ctx context.Context, wordDto *database_models.WordDto) *database_models.DatabaseError {
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
		return database_models.ProcessErrorFromDatabase(err, "SaveWords:QueryRowContext")
	}

	return nil
}

func (s *WordStorage) UpdateWord(ctx context.Context, wordDto *database_models.WordDto) *database_models.DatabaseError {
	// Implement updating logic here
	query := `UPDATE words SET
		word = $1,
		transcription = $2,
		meaning = $3,
		example = $4,
		word_level = $5,
		translation = $6
		WHERE id = $7
	`
	ctx, cancel := context.WithTimeout(ctx, QueryRowTimeout)
	defer cancel()

	_, err := s.db.ExecContext(
		ctx,
		query,
		wordDto.Word,
		wordDto.Transcription.String,
		wordDto.Meaning.String,
		wordDto.Example.String,
		wordDto.WordLevel.String,
		wordDto.Translation.String,
		wordDto.ID,
	)

	if err != nil {
		return database_models.ProcessErrorFromDatabase(err, "UpdateWords:ExecContext")
	}

	return nil
}
