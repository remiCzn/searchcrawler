package main

import "github.com/joho/godotenv"

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		panic(err)
	}

	crawler := initCrawler()

	for {
		crawler.step()
	}

}
