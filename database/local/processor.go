package database_local

import (
	"bytes"
	"errors"
	"strings"
	"unicode"

	"github.com/ledongthuc/pdf"
)

func GetContentFromPdf(fileName string) (map[string]int, int, error) {
	f, r, err := pdf.Open(locationForFile + fileName)
	m := make(map[string]int, 0)

	if err != nil {
		return m, -1, err
	}

	defer f.Close()

	var buf bytes.Buffer
	b, err := r.GetPlainText()

	if err != nil {
		return m, -1, err
	}

	buf.ReadFrom(b)

	all_words := strings.Split(buf.String(), " ")
	m, wordCount := getMapFromString(all_words)
	return m, wordCount, nil
}

func processExtraChar(st string) (string, error) {
	var buf bytes.Buffer
	var prev_symb = false
	for _, char := range st {
		if (char > 'a' && char <= 'z') || (char > 'A' && char <= 'Z') {
			prev_symb = false
			buf.WriteRune(unicode.ToLower(char))
		} else if char == ',' || char == '.' || char == '!' || char == '?' || char == ' ' || char == ':' || char == ';' || char == '(' || char == ')' || char == '\n' || char == '"' || char == '`' {
			if prev_symb {
				return "", errors.New("too much symbols " + st)
			}
			prev_symb = true
			continue
		} else if char == '-' || char == '\'' {
			if prev_symb {
				return "", errors.New("too much symbols " + st)
			}
			prev_symb = true
			buf.WriteRune(char)
		} else {
			return "", errors.New("invalid string " + st)
		}
	}
	if buf.Len() > 45 {
		return "", errors.New("string too long " + st)
	}
	return buf.String(), nil
}

func getMapFromString(all_words []string) (map[string]int, int) {
	m := make(map[string]int)
	wordCount := 0
	for _, word := range all_words {
		wd, err := processExtraChar(word)
		if err != nil || wd == "" {
			continue
		} else {
			wordCount++
			i, ok := m[wd]
			if ok {
				m[wd] = i + 1
			} else {
				m[wd] = 1
			}
		}
	}
	return m, wordCount
}
