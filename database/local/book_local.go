package database_local

import (
	"io"
	"mime/multipart"
	"os"
)

var locationForFile = "./database/local/"

func SaveBookToLocal(fileName string, file multipart.File) error {
	f, err := os.OpenFile(locationForFile+fileName, os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		return err
	}

	defer f.Close()
	io.Copy(f, file)
	return nil
}

func DeleteBookFromLocal(fileName string) error {
	return os.Remove(locationForFile + fileName)
}
