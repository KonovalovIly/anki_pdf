package model

import (
	"context"
	"database/sql"
	"io"
	"mime/multipart"
	"os"

	"github.com/KonovalovIly/anki_pdf/api/processor"
)

type BookDto struct {
	ID        int64          `json:"id"`
	Title     string         `json:"title"`
	AddedAt   string         `json:"added_at"`
	WordCount int            `json:"word_count"`
	WordMap   map[string]int `json:"-"`
}

type BookStorage struct {
	db *sql.DB
}

func (s *BookStorage) GetBook(ctx context.Context, bookID int64) (*BookDto, *DatabaseError) {
	// Implement getting logic here
	query := `SELECT id, title, added_at, words_count FROM books WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryRowTimeout)
	defer cancel()

	bookDto := &BookDto{}
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
		return nil, ProcessErrorFromDatabase(err)
	}

	return bookDto, nil
}

func (s *BookStorage) SaveBook(context context.Context, title string, file multipart.File, fileName string) (*BookDto, *DatabaseError) {
	// Implement saving logic here
	resultFileName := "./database/local/" + fileName
	err := saveBookToLocal(resultFileName, file)

	if err != nil {
		return nil, ProcessErrorFromDatabase(err)
	}

	bookDto, err := s.saveBookToDatabase(context, title)

	if err != nil {
		return nil, ProcessErrorFromDatabase(err)
	}

	// Implement processing logic here
	content, err := processor.GetContentFromPdf(resultFileName)
	if err != nil {
		return nil, ProcessErrorFromDatabase(err)
	}

	bookDto.WordMap, bookDto.WordCount = processor.ProcessContent(content)

	bookDto, e := s.UpdateBook(context, bookDto)

	if e != nil {
		return nil, e
	}

	err = deleteBookFromLocal(resultFileName)
	if err != nil {
		return nil, ProcessErrorFromDatabase(err)
	}

	return bookDto, nil
}

func (s *BookStorage) saveBookToDatabase(ctx context.Context, title string) (*BookDto, error) {
	query := `INSERT INTO books (title, user_id) VALUES ($1, $2) RETURNING id, added_at`
	ctx, cancel := context.WithTimeout(ctx, QueryRowTimeout)
	defer cancel()
	bookDto := &BookDto{
		Title: title,
	}

	err := s.db.QueryRowContext(
		ctx,
		query,
		title,
		1,
	).Scan(
		&bookDto.ID,
		&bookDto.AddedAt,
	)

	if err != nil {
		return nil, err
	}

	return bookDto, nil
}

func (s *BookStorage) UpdateBook(ctx context.Context, book *BookDto) (*BookDto, *DatabaseError) {
	query := `UPDATE books SET title = $1, words_count = $2 WHERE id = $3`

	_, err := s.db.ExecContext(
		ctx,
		query,
		book.Title,
		book.WordCount,
		book.ID,
	)
	if err != nil {
		return nil, ProcessErrorFromDatabase(err)
	}
	return book, nil
}

func saveBookToLocal(resultFileName string, file multipart.File) error {
	f, err := os.OpenFile(resultFileName, os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		return err
	}

	defer f.Close()
	io.Copy(f, file)
	return nil
}

func deleteBookFromLocal(resultFileName string) error {
	return os.Remove(resultFileName)
}
