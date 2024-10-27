package main

import (
	"fmt"
	"net/url"

	"github.com/gocolly/colly"
)

const INIT_PATH string = "https://fr.wikipedia.org/"

func main() {
	var visited = map[string]bool{}
	var waitingList = make([]string, 0)
	s := INIT_PATH

	for i := 0; i < 5; i++ {
		fmt.Printf("%v %v\n", len(waitingList), len(visited))
		waitingList, visited = crawl(s, waitingList, visited)
		s, waitingList = waitingList[0], waitingList[:1]
	}
}

func crawl(site string, waitingList []string, visited map[string]bool) ([]string, map[string]bool) {
	c := colly.NewCollector(
		colly.Async(true),
	)
	u, _ := url.Parse(site)
	baseUri := fmt.Sprintf("%s://%s", u.Scheme, u.Host)

	// called before an HTTP request is triggered
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL)
	})

	// triggered when the scraper encounters an error
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong: ", err)
	})

	// fired when the server responds
	c.OnResponse(func(r *colly.Response) {
		// fmt.Println("Page visited: ", r.Request.URL)
	})

	// triggered when a CSS selector matches an element
	c.OnHTML("a", func(e *colly.HTMLElement) {

		nextLink := e.Attr("href")

		u, _ := url.Parse(nextLink)

		if u.Host == "" {
			nextLink = fmt.Sprintf("%s%s", baseUri, nextLink)
		}

		_, ok := visited[nextLink]

		// println(nextLink)

		if !ok {
			waitingList = append(waitingList, nextLink)
			visited[nextLink] = true
		}
	})

	// triggered once scraping is done (e.g., write the data to a CSV file)
	c.OnScraped(func(r *colly.Response) {
		// fmt.Println(r.Request.URL, " scraped!")
	})

	c.Visit(site)
	c.Wait()

	return waitingList, visited
}
