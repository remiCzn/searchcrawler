package docparser

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func CreateDoc(site string) (*goquery.Document, error) {
	res, err := http.Get(site)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("status code error: %s | %d %s", site, res.StatusCode, res.Status))
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}
