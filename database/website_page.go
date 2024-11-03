package database

import (
	"fmt"
)

func (db *Database) PageExists(website_url string, website_page string) bool {
	exists, err := db.conn.Query(`SELECT 1 FROM website_page wp 
                                    JOIN website w ON w.id = wp.website_id
                                    where w.base_url = $1 and wp."path" = $2;`, website_url, website_page)

	CheckDbError(err, "PageExists:", website_url, website_page)

	defer exists.Close()

	return exists.Next()
}

func (db *Database) AddPage(website_url string, website_page string) error {
	_, err := db.conn.Exec(`INSERT INTO website_page (website_id, path, visited) 
                            VALUES ((SELECT id FROM website WHERE base_url = $1), $2, false);`, website_url, website_page)

	if err != nil {
		return err
	}
	return nil
}

type CrawlableWebsitePage struct {
	Id      int
	FullUrl string
	Visited bool
}

func (db *Database) GetNextPageToVisit() *CrawlableWebsitePage {
	var page CrawlableWebsitePage

	row, err := db.conn.Query(`SELECT wp.id, CONCAT(w.base_url, wp."path"), wp.visited 
                                FROM website_page wp
                                JOIN website w ON w.id = wp.website_id
                                WHERE 
                                    NOT wp.visited 
                                    AND (w.visited_last_time IS NULL OR w.visited_last_time <  now() - interval '10 second')
                                limit 1`)

	if err != nil {
		return nil
	}

	defer row.Close()

	if row.Next() {
		err = row.Scan(&page.Id, &page.FullUrl, &page.Visited)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return nil
		}
	}

	return &page
}

func (db *Database) SetPageVisited(id int) error {
	_, err := db.conn.Exec(`UPDATE website_page SET visited = true WHERE id = $1`, id)
	return err
}
