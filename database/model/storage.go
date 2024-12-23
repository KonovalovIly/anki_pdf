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
		NewWordsUser(ctx context.Context, userID int64, bookID int64, count int) ([]*WordDto, *DatabaseError)
	}

	Word interface {
		GetWord(ctx context.Context, text string) (*WordDto, *DatabaseError)
		GetWordById(ctx context.Context, wordId int64) (*WordDto, *DatabaseError)
		UpdateWords(ctx context.Context, wordDto *WordDto) (*WordDto, *DatabaseError)
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
