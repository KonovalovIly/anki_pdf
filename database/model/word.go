package database_models

import (
	"database/sql"
)

type WordDto struct {
	ID            int64          `json:"id"`
	Word          string         `json:"word"`
	Transcription sql.NullString `json:"transcription"`
	Meaning       sql.NullString `json:"meaning"`
	Example       sql.NullString `json:"example"`
	WordLevel     sql.NullString `json:"word_level"`
	Translation  sql.NullString `json:"translation"`
	Frequency     sql.NullInt16  `json:"frequency"`
}
