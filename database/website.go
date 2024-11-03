package database

import (
	"errors"
	docparser "searchcrawler/doc_parser"
	"time"
)

type Website struct {
	Id              int
	BaseUrl         string
	VisitedLastTime time.Time
	CreatedAt       time.Time
	Lang            string
}

func (db *Database) AddWebsite(base_url string) error {
	doc, err := docparser.CreateDoc(base_url)
	if err != nil {
		return err
	}
	lang, exists := doc.Find("html[lang]").First().Attr("lang")

	if !exists || len(lang) > 15 {
		return errors.New("invalid lang")
	}

	_, e := db.conn.Exec(`INSERT INTO website (base_url, lang) VALUES ($1, $2);`, base_url, lang)
	CheckDbError(e, "AddWebsite:", base_url, lang)
	return nil
}

func (db *Database) ExistsWebsite(base_url string) bool {
	rows, err := db.conn.Query(`SELECT 1 FROM website WHERE base_url = $1;`, base_url)
	CheckDbError(err, "ExistsWebsite:", base_url)

	defer rows.Close()

	return rows.Next()
}

func (db *Database) GetWebsite(full_url string) *Website {
	row, err := db.conn.Query(`SELECT id, base_url, lang FROM website WHERE base_url = $1;`, full_url)
	CheckDbError(err, "GetWebsite:", full_url)

	defer row.Close()

	if row.Next() {
		var website Website
		row.Scan(&website.Id, &website.BaseUrl, &website.Lang)
		return &website
	}
	return nil
}
