package database_models

type BookDto struct {
	ID        int64          `json:"id"`
	Title     string         `json:"title"`
	AddedAt   string         `json:"added_at"`
	WordCount int            `json:"word_count"`
	WordMap   map[string]int `json:"-"`
}
