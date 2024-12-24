package repository

import (
	"fmt"
	"strings"

	database_models "github.com/KonovalovIly/anki_pdf/database/model"
	"github.com/gocolly/colly"
)

func GetWordDetail(wordDto *database_models.WordDto) {
	c := colly.NewCollector(
		colly.AllowedDomains("dictionary.cambridge.org"),
	)
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"


	// get english level info
	c.OnHTML(".epp-xref.dxref", func(e *colly.HTMLElement) {
		if e.Index == 0 {
			wordDto.WordLevel.String = strings.TrimSpace(e.Text)
		}
	})

	// get meanings
	var meaningSelector = ".def.ddef_d.db"

	c.OnHTML(meaningSelector, func(e *colly.HTMLElement) {
		result := strings.TrimSpace(e.Text)
		result = strings.Replace(result, ":", "", -1)
		wordDto.Meaning.String = result
	})

	// get example sentence
	c.OnHTML(".eg.deg", func(e *colly.HTMLElement) {
		if e.Index == 0 {
			result := e.Text
			wordDto.Example.String = strings.TrimSpace(result)
		}
	})

	c.Visit(fmt.Sprintf("https://dictionary.cambridge.org/dictionary/english/%v", wordDto.Word))
}
