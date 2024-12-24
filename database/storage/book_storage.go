package storage

import (
	"context"
	"database/sql"

	database_models "github.com/KonovalovIly/anki_pdf/database/model"
)

type BookStorage struct {
	db *sql.DB
}

func (s *BookStorage) GetBook(ctx context.Context, bookID int64) (*database_models.BookDto, *database_models.DatabaseError) {
	query := `SELECT
		id,
		title,
		added_at,
		words_count
		FROM books WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryRowTimeout)
	defer cancel()

	bookDto := &database_models.BookDto{}
	err := s.db.QueryRowContext(
		ctx,
		query,
		bookID,
	).Scan(
		&bookDto.ID,
		&bookDto.Title,
		&bookDto.AddedAt,
		&bookDto.WordCount,
	)

	if err != nil {
		return nil, database_models.ProcessErrorFromDatabase(err, "GetBook:QueryRowContext")
	}

	return bookDto, nil
}

func (s *BookStorage) SaveBook(ctx context.Context, book *database_models.BookDto, userId int64) *database_models.DatabaseError {
	query := `INSERT INTO books
		(title, user_id)
		VALUES ($1, $2) RETURNING id, added_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryRowTimeout)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		book.Title,
		userId,
	).Scan(
		&book.ID,
		&book.AddedAt,
	)

	if err != nil {
		return database_models.ProcessErrorFromDatabase(err, "SaveBook:QueryRowContext")
	}

	return nil
}

func (s *BookStorage) UpdateBook(ctx context.Context, book *database_models.BookDto) *database_models.DatabaseError {
	query := `UPDATE books SET
		title = $1,
		words_count = $2
		WHERE id = $3
	`

	_, err := s.db.ExecContext(
		ctx,
		query,
		book.Title,
		book.WordCount,
		book.ID,
	)
	if err != nil {
		return database_models.ProcessErrorFromDatabase(err, "UpdateBook:ExecContext")
	}
	return nil
}
