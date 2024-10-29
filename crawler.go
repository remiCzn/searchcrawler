package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

type Crawler struct {
	visited       map[string]bool
	waitingList   []string
	site          string
	robotsChecker *RobotChecker
}

func initCrawler() *Crawler {
	c := Crawler{}
	c.waitingList = []string{"https://membre.leadersante-groupe.fr/"}
	c.visited = map[string]bool{}
	c.robotsChecker = &RobotChecker{}
	c.robotsChecker.init()
	for _, el := range c.waitingList {
		c.visited[el] = true
	}

	return &c
}

func (c *Crawler) step() {
	c.site = c.waitingList[0]
	c.waitingList = c.waitingList[1:]

	max := len(c.waitingList)
	if max > 5 {
		max = 5
	}

	c.crawl()
}

func (c *Crawler) crawl() {
	fmt.Println("Crawling:", c.site)
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

	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Attr("href")

		if exists {
			fmt.Println(link)

			u, err := url.Parse(link)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			if u.Host == "" {
				link = fmt.Sprintf("%s%s", baseUri, link)
			}
			_, ok := c.visited[link]
			if !ok {
				c.visited[link] = true

				if c.robotsChecker.checkIfAllowed(link) {
					c.waitingList = append(c.waitingList, link)
				}
			}
		}
	})
}

func (c *Crawler) printStats() {
	fmt.Print("Sites checked:", len(c.visited), "\n Remaining sites:", len(c.waitingList), "\n")
}
