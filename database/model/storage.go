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
		SaveBook(ctx context.Context, title string, book multipart.File, fileName string) (*BookDto, error)
	}
	Word interface {
		SaveWords(ctx context.Context, book *BookDto, wordDto *WordDto) (*WordDto, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Book: &BookStorage{db: db},
		Word: &WordStorage{db: db},
	}
}
