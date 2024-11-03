package main

import (
	"fmt"
	"net/url"
	"searchcrawler/database"
	docparser "searchcrawler/doc_parser"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Crawler struct {
	site          string
	robotsChecker *RobotChecker
	db            *database.Database
}

func initCrawler() *Crawler {
	c := Crawler{}
	c.robotsChecker = &RobotChecker{}
	c.robotsChecker.init()

	c.db = &database.Database{}
	c.db.Init()

	return &c
}

func (c *Crawler) step() {
	site := c.db.GetNextPageToVisit()

	if site == nil {
		fmt.Println("No more website to crawl (OoO)")
		time.Sleep(10 * time.Second)
	}
	c.site = site.FullUrl
	c.crawl()

	err := c.db.SetPageVisited(site.Id)
	if err != nil {
		fmt.Println("Error setting page visited:", err)
	}
}

func (c *Crawler) crawl() {
	fmt.Println("Crawling:", c.site)
	doc, err := docparser.CreateDoc(c.site)
	if err != nil {
		fmt.Println("Error creating doc:", err)
		return
	}

	u, _ := url.Parse(c.site)
	baseUri := fmt.Sprintf("%s://%s", u.Scheme, u.Host)

	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Attr("href")

		if exists {
			u, err := url.Parse(link)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			if u.Host == "" {
				link = fmt.Sprintf("%s%s", baseUri, link)
			}

			if c.robotsChecker.checkIfAllowed(link) {
				c.addUriToWaitingList(link)
			}
		}
	})
}

func (c *Crawler) addUriToWaitingList(link string) {

	u, _ := url.Parse(link)
	baseUrl := fmt.Sprintf("%s://%s", u.Scheme, u.Host)
	path := u.Path
	if u.RawQuery != "" {
		path = fmt.Sprintf("%s?%s", u.Path, u.RawQuery)
	}
	if u.Fragment != "" {
		path = fmt.Sprintf("%s#%s", u.Path, u.Fragment)
	}

	// Add root website to db
	if !c.db.ExistsWebsite(baseUrl) {
		err := c.db.AddWebsite(baseUrl)
		if err != nil {
			fmt.Println("Error adding website to db:", err)
			return
		}
	}

	// Check if the site is in english or french
	website := c.db.GetWebsite(baseUrl)
	if website == nil || (website.Lang != "en" && website.Lang != "fr") {
		return
	}

	// Check if the website respect the robots.txt
	if !c.robotsChecker.checkIfAllowed(link) {
		return
	}

	//Check if the page is already visited
	if c.db.PageExists(baseUrl, path) {
		return
	}

	fmt.Println(fmt.Sprintf("%s%s", baseUrl, path))

	err := c.db.AddPage(baseUrl, path)

	if err != nil {
		fmt.Println(err)
	}
}
