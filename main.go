package main

func main() {
	crawler := initCrawler()

	for i := 0; i < 20; i++ {
		crawler.step()
	}

	crawler.printStats()

}
