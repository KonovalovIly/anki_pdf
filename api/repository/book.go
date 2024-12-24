package repository

import (
	"context"
	"mime/multipart"

	database_local "github.com/KonovalovIly/anki_pdf/database/local"
	database_models "github.com/KonovalovIly/anki_pdf/database/model"
	storage "github.com/KonovalovIly/anki_pdf/database/storage"
)

func ProcessUploadBook(
	context context.Context,
	book *database_models.BookDto,
	file multipart.File,
	fileName string,
	storage *storage.Storage,
	userId int64,
) *database_models.DatabaseError {

	err := SaveBook(context, book, file, fileName, storage, userId)
	if err != nil {
		return err
	}

	wordsMap := book.WordMap

	for word := range wordsMap {
		wordDto, err := storage.Word.GetWordByName(context, word)

		if err != nil && err.Typ == "no_row" {

			wordDto = &database_models.WordDto{}
			wordDto.Word = word

			err := storage.Word.SaveWord(context, wordDto)
			if err != nil {
				return err
			}

		} else if err != nil {
			return err
		}

		err = storage.BookWord.SaveWordWithBookConnection(context, book, wordDto)
		if err != nil {
			return err
		}
	}

	return nil
}

func SaveBook(
	context context.Context,
	book *database_models.BookDto,
	file multipart.File,
	fileName string,
	storage *storage.Storage,
	userId int64,
) *database_models.DatabaseError {

	err := database_local.SaveBookToLocal(fileName, file)

	if err != nil {
		return database_models.ProcessErrorFromDatabase(err, "SaveBook:SaveBookToLocal")
	}

	exception := storage.Book.SaveBook(context, book, userId)

	if exception != nil {
		return exception
	}

	book.WordMap, book.WordCount, err = database_local.GetContentFromPdf(fileName)

	if err != nil {
		return database_models.ProcessErrorFromDatabase(err, "SaveBook:GetContentFromPdf")
	}

	exception = storage.Book.UpdateBook(context, book)

	if exception != nil {
		return exception
	}

	err = database_local.DeleteBookFromLocal(fileName)
	if err != nil {
		return database_models.ProcessErrorFromDatabase(err, "SaveBook:DeleteBookFromLocal")
	}

	return nil
}
