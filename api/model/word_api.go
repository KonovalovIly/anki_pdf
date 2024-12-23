package api_model

import "github.com/KonovalovIly/anki_pdf/database/model"

type WordApi struct {
	ID            int64  `json:"id"`
	Word          string `json:"word"`
	Transcription string `json:"transcription"`
	Meaning       string `json:"meaning"`
	Example       string `json:"example"`
	WordLevel     string `json:"word_level"`
	Translations  string `json:"translations"`
	Frequency     int    `json:"frequency"`
}

func MapListDtoToApiWord(dto []*model.WordDto) []*WordApi {
	apiWords := make([]*WordApi, 0, len(dto))
	for _, word := range dto {
		apiWords = append(apiWords, MapDtoToApiWord(word))
	}
	return apiWords
}

func MapDtoToApiWord(dto *model.WordDto) *WordApi {
	return &WordApi{
		ID:            dto.ID,
		Word:          dto.Word,
		Transcription: dto.Transcription.String,
		Meaning:       dto.Meaning.String,
		Example:       dto.Example.String,
		WordLevel:     dto.WordLevel.String,
		Translations:  dto.Translations.String,
		Frequency:     int(dto.Frequency.Int16),
	}
}
