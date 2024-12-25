package route

import (
	"fmt"
	"net/http"

	api_model "github.com/KonovalovIly/anki_pdf/api/model"
	"github.com/KonovalovIly/anki_pdf/api/repository"
	api_utils "github.com/KonovalovIly/anki_pdf/api/utils"
	"github.com/go-chi/chi/v5"
)

func (app *Application) WordsHandlersSetup(r chi.Router) {
	r.Get("/new_words", app.getNewWordsForBookHandler)
}

func (app *Application) getNewWordsForBookHandler(w http.ResponseWriter, r *http.Request) {
	var payload api_model.NewWordsPayload

	if err := api_utils.ReadJson(w, r, &payload); err != nil {
		api_utils.WriteJsonError(w, http.StatusBadRequest, err)
		return
	}

	if err := api_utils.Validator.Struct(payload); err != nil {
		api_utils.WriteJsonError(w, http.StatusBadRequest, err)
		return
	}

	bookId := payload.BookId
	count := payload.Count

	if bookId <= 0 {
		api_utils.WriteJsonError(w, http.StatusBadRequest, fmt.Errorf("invalid book ID"))
		return
	}

	ctx := r.Context()
	_, e := app.Storage.Book.GetBook(ctx, bookId)
	if e != nil {
		api_utils.WriteJsonDatabaseError(w, http.StatusInternalServerError, e)
		return
	}

	words, e := app.Storage.UserWord.NewWordsUser(ctx, 1, bookId, count)

	for i := range words {
		currentWord := words[i]
		wordDto, e := app.Storage.Word.GetWordById(ctx, currentWord.ID)
		if e != nil {
			api_utils.WriteJsonDatabaseError(w, http.StatusInternalServerError, e)
			return
		}

		if !wordDto.Meaning.Valid {
			repository.GetWordDetail(wordDto)
			_ = app.Storage.Word.UpdateWord(ctx, wordDto)
		}

		words[i].Word = wordDto.Word
		words[i].Meaning = wordDto.Meaning
		words[i].Example = wordDto.Example
		words[i].WordLevel = wordDto.WordLevel
	}

	if e != nil {
		api_utils.WriteJsonDatabaseError(w, http.StatusInternalServerError, e)
		return
	}

	api_utils.JsonResponse(w, http.StatusOK, api_model.MapListDtoToApiWord(words))
}