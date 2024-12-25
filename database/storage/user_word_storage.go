package storage

import (
	"context"
	"database/sql"

	database_models "github.com/KonovalovIly/anki_pdf/database/model"
)

type UserWordStorage struct {
	db *sql.DB
}

func (s *UserWordStorage) NewWordsUser(ctx context.Context, userID int64, bookID int64, count int) ([]*database_models.WordDto, *database_models.DatabaseError) {
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
		return nil, database_models.ProcessErrorFromDatabase(err, "NewWordsUser:QueryContext")
	}

	defer rows.Close()

	wordDtos := make([]*database_models.WordDto, count)
	i := 0

	for rows.Next() {
		wordDto := &database_models.WordDto{}
		err := rows.Scan(
			&wordDto.ID,
			&wordDto.Frequency,
		)
		if err != nil {
			return nil, database_models.ProcessErrorFromDatabase(err, "NewWordsUser:Scan")
		}
		wordDtos[i] = wordDto
		i++
	}

	return wordDtos, nil
}

func (s *UserWordStorage) KnownWordsBook(ctx context.Context, userID int64, bookID int64) (*database_models.BookWithNounWords, *database_models.DatabaseError) {
	query := `
        SELECT b.book_id, SUM (b.frequency), bb.title
		FROM books_words b
		JOIN users_words u ON b.word_id=u.word_id
		JOIN books bb ON b.book_id = bb.id
		WHERE u.user_id = $1 AND u.is_learned = TRUE AND b.book_id = $2
		GROUP BY b.book_id, bb.title;
    `

	ctx, cancel := context.WithTimeout(ctx, QueryRowTimeout)
	defer cancel()

	bookWithNounWords := &database_models.BookWithNounWords{}
	err := s.db.QueryRowContext(
		ctx,
		query,
		userID,
		bookID,
	).Scan(
		&bookWithNounWords.ID,
		&bookWithNounWords.AlreadyKnownWords,
		&bookWithNounWords.Title,
	)

	if err != nil {
		return nil, database_models.ProcessErrorFromDatabase(err, "KnownWordsBook:QueryRowContext")
	}
	return bookWithNounWords, nil
}
