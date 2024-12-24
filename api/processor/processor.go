package processor

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/ledongthuc/pdf"
)

func GetContentFromPdf(path string) (string, error) {
	f, r, err := pdf.Open(path)
	// remember close file
	defer f.Close()
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	buf.ReadFrom(b)
	return buf.String(), nil
}

func ProcessContent(content string) (map[string]int, int) {
	fmt.Println("Processing content...")
	all_words := strings.Split(content, " ")
	fmt.Println("Content Splited...")
	m, wordCount := getMapFromString(all_words)
	fmt.Println("Content Processed...")
	return m, wordCount
}

func processExtraChar(st string) (string, error) {
	var buf bytes.Buffer
	var prev_symb = false
	for _, char := range st {
		if unicode.IsLetter(char) {
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
