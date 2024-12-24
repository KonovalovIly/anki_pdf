package api_model

type NewWordsPayload struct {
	BookId int64 `json:"book_id"  validate:"required"`
	Count  int   `json:"count"`
}
