package storage

import (
	"context"
	"database/sql"

	database_models "github.com/KonovalovIly/anki_pdf/database/model"
)

type BookWordStorage struct {
	db *sql.DB
}

func (s *BookWordStorage) SaveWordWithBookConnection(ctx context.Context, book *database_models.BookDto, word *database_models.WordDto) *database_models.DatabaseError {
	query := `
		INSERT INTO books_words
		(book_id, word_id, frequency)
		VALUES ($1, $2, $3)
	`

	_, err := s.db.ExecContext(
		ctx,
		query,
		book.ID,
		word.ID,
		book.WordMap[word.Word],
	)

	if err != nil {
		return database_models.ProcessErrorFromDatabase(err, "SaveWordWithBookConnection:ExecContext")
	}

	return nil
}
