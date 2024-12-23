package main

import (
	"fmt"
	"net/http"
	"strconv"

	api_model "github.com/KonovalovIly/anki_pdf/api/model"
	"github.com/go-chi/chi/v5"
)

func (app *Application) getNewWordsForBookHandler(w http.ResponseWriter, r *http.Request) {
	bookID, err := strconv.ParseInt(chi.URLParam(r, "bookID"), 10, 64)
	if err != nil || bookID <= 0 {
		app.writeJsonError(w, http.StatusBadRequest, fmt.Errorf("Invalid book ID"))
		return
	}
	ctx := r.Context()

	words, e := app.Storage.Book.NewWordsUser(ctx, 1, bookID, 30)
	if e != nil {
		app.writeJsonDatabaseError(w, http.StatusInternalServerError, e)
		return
	}

	app.jsonResponse(w, http.StatusOK, api_model.MapListDtoToApiWord(words))
}
