package model

import (
	"context"
	"database/sql"
	"mime/multipart"
	"time"
)

var QueryRowTimeout = time.Second * 5

type Storage struct {
	Book interface {
		GetBook(ctx context.Context, bookID int64) (*BookDto, *DatabaseError)
		SaveBook(ctx context.Context, title string, book multipart.File, fileName string) (*BookDto, *DatabaseError)
	}
	Word interface {
		GetWord(ctx context.Context, text string) (*WordDto, *DatabaseError)
		SaveWords(ctx context.Context, wordDto *WordDto) (*WordDto, *DatabaseError)
		SaveWordWithBookConnection(ctx context.Context, book *BookDto, wordDto *WordDto) *DatabaseError
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Book: &BookStorage{db: db},
		Word: &WordStorage{db: db},
	}
}
