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
		return nil, ProcessErrorFromDatabase(err, "GetBook:44")
	}

	return bookDto, nil
}

func (s *BookStorage) SaveBook(context context.Context, title string, file multipart.File, fileName string) (*BookDto, *DatabaseError) {
	// Implement saving logic here
	resultFileName := "./database/local/" + fileName
	err := saveBookToLocal(resultFileName, file)

	if err != nil {
		return nil, ProcessErrorFromDatabase(err, "SaveBook:56")
	}

	bookDto, err := s.saveBookToDatabase(context, title)

	if err != nil {
		return nil, ProcessErrorFromDatabase(err, "SaveBook:62")
	}

	// Implement processing logic here
	content, err := processor.GetContentFromPdf(resultFileName)
	if err != nil {
		return nil, ProcessErrorFromDatabase(err, "SaveBook:68")
	}

	bookDto.WordMap, bookDto.WordCount = processor.ProcessContent(content)

	bookDto, e := s.UpdateBook(context, bookDto)

	if e != nil {
		return nil, e
	}

	err = deleteBookFromLocal(resultFileName)
	if err != nil {
		return nil, ProcessErrorFromDatabase(err, "SaveBook:81")
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
		return nil, ProcessErrorFromDatabase(err, "UpdateBook:123")
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

func (s *BookStorage) NewWordsUser(ctx context.Context, userID int64, bookID int64, count int) ([]*WordDto, *DatabaseError) {
	query := `
		SELECT bw.word_id, SUM(bw.frequency) AS total_frequency
		FROM books_words bw
		LEFT JOIN users_words uw ON bw.word_id = uw.word_id AND uw.user_id = $1
		WHERE bw.book_id = $2
		AND (uw.user_id IS NULL OR uw.is_learned = FALSE)
		GROUP BY bw.word_id
		ORDER BY total_frequency DESC
		LIMIT $3;
	`

	rows, err := s.db.QueryContext(
		ctx,
		query,
		userID,
		bookID,
		count,
	)

	if err != nil {
		return nil, ProcessErrorFromDatabase(err, "NewWordsUser:165")
	}
	defer rows.Close()
	wordDtos := make([]*WordDto, 0)
	for rows.Next() {
		wordDto := &WordDto{}
		err := rows.Scan(
			&wordDto.ID,
			&wordDto.Frequency,
		)
		if err != nil {
			return nil, ProcessErrorFromDatabase(err, "NewWordsUser:176")
		}
		wordDtos = append(wordDtos, wordDto)
	}
	return wordDtos, nil
}
