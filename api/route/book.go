package route

import (
	"fmt"
	"net/http"
	"strconv"

	repository "github.com/KonovalovIly/anki_pdf/api/repository"
	api_utils "github.com/KonovalovIly/anki_pdf/api/utils"
	database_models "github.com/KonovalovIly/anki_pdf/database/model"
	"github.com/go-chi/chi/v5"
)

func (app *Application) BookHandlerSetup(r chi.Router) {
	r.Route("/book", func(r chi.Router) {
		r.Post("/upload", app.bookUploadHandler)

		r.Route("/{bookID}", func(r chi.Router) {
			r.Get("/", app.bookGetHandler)
			r.Get("/known_words", app.knownWordsHandler)
		})
	})
}

func (app *Application) bookGetHandler(w http.ResponseWriter, r *http.Request) {
	bookID, err := strconv.ParseInt(chi.URLParam(r, "bookID"), 10, 64)
	if err != nil || bookID <= 0 {
		api_utils.WriteJsonError(w, http.StatusBadRequest, fmt.Errorf("invalid book ID"))
		return
	}

	bookDto, e := app.Storage.Book.GetBook(r.Context(), bookID)
	if e != nil {
		api_utils.WriteJsonDatabaseError(w, http.StatusInternalServerError, e)
		return
	}

	api_utils.JsonResponse(w, http.StatusAccepted, bookDto)
}

func (app *Application) bookUploadHandler(w http.ResponseWriter, r *http.Request) {
	bookTitle := r.Header.Get("Book_Title")
	bookLang := r.Header.Get("Book_Lang")
	if bookTitle == "" {
		api_utils.WriteJsonError(w, http.StatusBadRequest, fmt.Errorf("book title is required"))
		return
	}

	if bookLang == "" {
		api_utils.WriteJsonError(w, http.StatusBadRequest, fmt.Errorf("book lang is required"))
		return
	}

	file, fileHeader, err := r.FormFile("fileupload")
	if err != nil {
		api_utils.WriteJsonError(w, http.StatusBadRequest, fmt.Errorf("error reading book file"))
		return
	}

	defer file.Close()
	ctx := r.Context()

	bookDto := &database_models.BookDto{
		Title: bookTitle,
	}

	e := repository.ProcessUploadBook(ctx, bookDto, file, fileHeader.Filename, &app.Storage, 1)
	if e != nil {
		api_utils.WriteJsonDatabaseError(w, http.StatusBadRequest, e)
		return
	}

	api_utils.JsonResponse(w, http.StatusAccepted, bookDto)
}

func (app *Application) knownWordsHandler(w http.ResponseWriter, r *http.Request) {
	bookID, err := strconv.ParseInt(chi.URLParam(r, "bookID"), 10, 64)
	if err != nil || bookID <= 0 {
		api_utils.WriteJsonError(w, http.StatusBadRequest, fmt.Errorf("invalid book ID"))
		return
	}

	ctx := r.Context()
	words, e := app.Storage.UserWord.KnownWordsBook(ctx, 1, bookID)

	if e != nil {
		api_utils.WriteJsonDatabaseError(w, http.StatusInternalServerError, e)
		return
	}

	api_utils.JsonResponse(w, http.StatusOK, words)
}
