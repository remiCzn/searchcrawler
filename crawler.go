package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

type Crawler struct {
	visited     map[string]bool
	waitingList []string
	site        string
}

func initCrawler() *Crawler {
	c := Crawler{}
	c.waitingList = []string{"https://fr.wikipedia.org/", "https://www.lemonde.fr/"}
	c.visited = map[string]bool{}
	for _, el := range c.waitingList {
		c.visited[el] = true
	}

	return &c
}

func (c *Crawler) step() {
	c.site, c.waitingList = c.waitingList[0], c.waitingList[:1]
	c.crawl()
}

func (c *Crawler) crawl() {
	res, err := http.Get(c.site)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	u, _ := url.Parse(c.site)
	baseUri := fmt.Sprintf("%s://%s", u.Scheme, u.Host)

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Attr("href")

		if exists {
			u, _ := url.Parse(link)
			if u.Host == "" {
				link = fmt.Sprintf("%s%s", baseUri, link)
			}
			_, ok := c.visited[link]
			if !ok {
				fmt.Println(link)
				c.visited[link] = true
				c.waitingList = append(c.waitingList, link)
			}
		}
	})
}
