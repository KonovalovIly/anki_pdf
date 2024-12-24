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
