package main

import (
	"fmt"
	"net/http"
	"strconv"

	repository "github.com/KonovalovIly/anki_pdf/api/repository"
	database_models "github.com/KonovalovIly/anki_pdf/database/model"
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
	bookLang := r.Header.Get("Book_Lang")
	if bookTitle == "" {
		app.writeJsonError(w, http.StatusBadRequest, fmt.Errorf("Book title is required"))
		return
	}

	if bookLang == "" {
		app.writeJsonError(w, http.StatusBadRequest, fmt.Errorf("Book lang is required"))
		return
	}

	file, fileHeader, err := r.FormFile("fileupload")
	if err != nil {
		app.writeJsonError(w, http.StatusBadRequest, fmt.Errorf("Error reading book file"))
		return
	}

	defer file.Close()
	ctx := r.Context()

	bookDto := &database_models.BookDto{
		Title: bookTitle,
	}

	e := repository.ProcessUploadBook(ctx, bookDto, file, fileHeader.Filename, &app.Storage, 1)
	if err != nil {
		app.writeJsonDatabaseError(w, http.StatusBadRequest, e)
		return
	}

	app.jsonResponse(w, http.StatusAccepted, bookDto)
}
