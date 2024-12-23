package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/KonovalovIly/anki_pdf/database/model"
	"github.com/go-chi/chi/v5"
)

func (app *Application) bookGetHandler(w http.ResponseWriter, r *http.Request) {
	bookID, err := strconv.ParseInt(chi.URLParam(r, "bookID"), 10, 64)
	if err != nil || bookID <= 0 {
		app.writeJsonError(w, http.StatusBadRequest, fmt.Errorf("Invalid book ID"))
		return
	}

	bookDto, e := app.Storage.Book.GetBook(r.Context(), bookID)
	if e != nil {
		app.writeJsonDatabaseError(w, http.StatusInternalServerError, e)
		return
	}

	app.jsonResponse(w, http.StatusAccepted, bookDto)
}

func (app *Application) bookUploadHandler(w http.ResponseWriter, r *http.Request) {
	bookTitle := r.Header.Get("Book_Title")
	if bookTitle == "" {
		app.writeJsonError(w, http.StatusBadRequest, fmt.Errorf("Book title is required"))
		return
	}

	file, fileHeader, err := r.FormFile("fileupload")
	if err != nil {
		app.writeJsonError(w, http.StatusBadRequest, fmt.Errorf("Error reading book file"))
		return
	}

	defer file.Close()

	bookDto, erro := app.Storage.Book.SaveBook(r.Context(), bookTitle, file, fileHeader.Filename)
	if erro != nil {
		app.writeJsonDatabaseError(w, http.StatusInternalServerError, erro)
		return
	}

	wordsMap := bookDto.WordMap

	for word := range wordsMap {
		wordDto, err := app.Storage.Word.GetWord(r.Context(), word)

		if err != nil && err.Typ == "no_row" {
			wordDto = &model.WordDto{}
			wordDto.Word = word
			app.saveNewWord(r, wordDto)
		} else if err != nil {
			app.writeJsonDatabaseError(w, http.StatusInternalServerError, err)
			return
		}

		err = app.Storage.Word.SaveWordWithBookConnection(r.Context(), bookDto, wordDto)
		if err != nil {
			app.writeJsonDatabaseError(w, http.StatusInternalServerError, err)
			return
		}
	}

	app.jsonResponse(w, http.StatusAccepted, bookDto)
}

func (app *Application) saveNewWord(r *http.Request, wordDto *model.WordDto) {
	app.Storage.Word.SaveWords(r.Context(), wordDto)
}
