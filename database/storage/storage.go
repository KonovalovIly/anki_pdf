package storage

import (
	"context"
	"database/sql"
	"time"

	database_models "github.com/KonovalovIly/anki_pdf/database/model"
)

var QueryRowTimeout = time.Second * 5

type Storage struct {
	Book interface {
		GetBook(ctx context.Context, bookID int64) (*database_models.BookDto, *database_models.DatabaseError)
		SaveBook(ctx context.Context, book *database_models.BookDto, userId int64) *database_models.DatabaseError
		UpdateBook(ctx context.Context, book *database_models.BookDto) *database_models.DatabaseError
	}

	Word interface {
		GetWordById(ctx context.Context, wordId int64) (*database_models.WordDto, *database_models.DatabaseError)
		GetWordByName(ctx context.Context, text string) (*database_models.WordDto, *database_models.DatabaseError)
		SaveWord(ctx context.Context, wordDto *database_models.WordDto) *database_models.DatabaseError
		UpdateWord(ctx context.Context, wordDto *database_models.WordDto) *database_models.DatabaseError
	}

	BookWord interface {
		SaveWordWithBookConnection(ctx context.Context, book *database_models.BookDto, word *database_models.WordDto) *database_models.DatabaseError
	}

	UserWord interface {
		NewWordsUser(ctx context.Context, userID int64, bookID int64, count int) ([]*database_models.WordDto, *database_models.DatabaseError)
		KnownWordsBook(ctx context.Context, userID int64, bookID int64) (*database_models.BookWithNounWords, *database_models.DatabaseError)
		MarkAsLearned(ctx context.Context, userID int64, wordID int64) *database_models.DatabaseError
	}

	User interface {
		GetUser(ctx context.Context, userID int64) (*database_models.UserDto, *database_models.DatabaseError)
		SaveUser(ctx context.Context, userDto *database_models.UserDto) *database_models.DatabaseError
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Book:     &BookStorage{db: db},
		Word:     &WordStorage{db: db},
		BookWord: &BookWordStorage{db: db},
		UserWord: &UserWordStorage{db: db},
		User:     &UserStorage{db: db},
	}
}
