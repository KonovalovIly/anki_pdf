package route

import (
	"fmt"
	"net/http"

	api_model "github.com/KonovalovIly/anki_pdf/api/model"
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

	words, e := app.Storage.UserWord.NewWordsUser(ctx, 1, bookId, count)

	for i := range words {
		currentWord := words[i]
		wordDto, e := app.Storage.Word.GetWordById(ctx, currentWord.ID)
		if e != nil {
			api_utils.WriteJsonDatabaseError(w, http.StatusInternalServerError, e)
			return
		}
		currentWord.Word = wordDto.Word
	}

	if e != nil {
		api_utils.WriteJsonDatabaseError(w, http.StatusInternalServerError, e)
		return
	}

	api_utils.JsonResponse(w, http.StatusOK, api_model.MapListDtoToApiWord(words))
}
