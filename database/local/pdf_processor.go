package database_local

import (
	"errors"
	"image/jpeg"
	"os"
	"strings"

	"github.com/gen2brain/go-fitz"
	"github.com/google/uuid"
)

func GetPreviewImage(filename string) (string, error) {
	str, ok := strings.CutSuffix(filename, ".pdf")
	if !ok {
		return "", errors.New("file name must have a .pdf extension")
	}

	file := locationForFile + filename
	resultFile := locationForFile + uuid.New().String() + str + ".jpg"

	doc, err := fitz.New(file)
	if err != nil {
		return "", err
	}
	defer doc.Close()

	img, err := doc.Image(0)
	if err != nil {
		return "", err
	}

	f, err := os.Create(resultFile)
	if err != nil {
		return "", err
	}

	err = jpeg.Encode(f, img, &jpeg.Options{Quality: jpeg.DefaultQuality})
	if err != nil {
		return "", err
	}

	f.Close()
	return resultFile, nil
}

func DeletePreviewImage(filename string) error {
	return os.Remove(locationForFile + filename)
}
