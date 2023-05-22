package main

import (
	"crypto/tls"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/matinkhosravani/fidibo_crawler/app"
	"github.com/matinkhosravani/fidibo_crawler/utils"
	"log"
	"net/http"
	"strings"
	"sync"
)

func main() {
	app.LoadEnv()
	// Instantiate default collector
	c := colly.NewCollector(
		colly.AllowedDomains("fidibo.com"),
	)
	book := c.Clone()
	s := utils.SetupCollyStorage(book)
	// close redis client
	defer s.Client.Close()

	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	})

	var wg sync.WaitGroup
	c.OnXML("//sitemapindex/sitemap/loc", func(e *colly.XMLElement) {
		link := e.Text
		if strings.Contains(link, "sitemap_book") {
			book.Visit(link)
		}
	})

	book.OnXML("//urlset/url/loc", func(e *colly.XMLElement) {
		bookID, err := utils.ExtractBoodIDFromURL(e.Text)
		if err != nil {
			log.Println(err)
		} else {
			wg.Add(1)
			go func(bookID string) {
				defer wg.Done()
				s.Client.HSet("fidibo_books", bookID, e.Text)
				s.Client.SAdd("fidibo_book_ids", bookID)
			}(bookID)
		}
	})

	err := c.Visit("https://fidibo.com/sitemap_index.xml")
	if err != nil {
		fmt.Println(err)
	}

	wg.Wait()

	log.Println("Finished Finding Book id from sitemap")
}
