package main

import (
	"net/http"
	"os"

	"github.com/KonovalovIly/anki_pdf/database/model"
)

func (app *Application) bookUploadHandler(w http.ResponseWriter, r *http.Request) {
	bookTitle := r.Header.Get("Book_Title")
	if bookTitle == "" {
		app.writeJsonError(w, http.StatusBadRequest, "Book title is required")
		return
	}

	file, fileHeader, err := r.FormFile("fileupload")
	if err != nil {
		app.writeJsonError(w, http.StatusBadRequest, "Error reading book file")
		return
	}
	defer file.Close()

	bookDto, err := app.Storage.Book.SaveBook(r.Context(), bookTitle, file, fileHeader.Filename)
	if err != nil {
		app.writeJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	wordsMap := bookDto.WordMap

	for wd, _ := range wordsMap {
		word := model.WordDto{
			Word: wd,
		}
		_, err := app.Storage.Word.SaveWords(r.Context(), bookDto, &word)
		if err != nil {
			app.writeJsonError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	app.jsonResponse(w, http.StatusAccepted, bookDto)
}

func (app *Application) bookDeleteHandler(w http.ResponseWriter, r *http.Request) {
	err := os.Remove("./database/local/CSAPP_2016.pdf")
	if err != nil {
		app.writeJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	app.jsonResponse(w, http.StatusOK, nil)
}
