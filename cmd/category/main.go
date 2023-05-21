package main

import (
	"crypto/tls"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/matinkhosravani/fidibo_crawler/app"
	"github.com/matinkhosravani/fidibo_crawler/storage"
	"github.com/matinkhosravani/fidibo_crawler/utils"
	"net/http"
	"regexp"
	"strings"
)

func main() {
	app.LoadEnv()

	// Instantiate default collector
	c := colly.NewCollector(
		colly.AllowedDomains("fidibo.com"),
		colly.URLFilters(
			regexp.MustCompile(".*((fidibo\\.com/?$)|(fidibo\\.com/category/[^/]*/?(\\?page=\\d)?(&keyword=.*)?(&sorting=.*)?)$)"),
		),
		colly.MaxDepth(50),
		colly.Async(true),
	)
	s := storage.SetupCollyStorage(c)
	// close redis client
	defer s.Client.Close()

	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	})

	// On every an element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		link = strings.Replace(link, "?&", "?", 1)
		// Visit link found on page
		// Only those links are visited which are matched by  any of the URLFilter regexps
		if strings.Contains(link, "/book/") {
			url := e.Request.AbsoluteURL(link)
			bookID, err := utils.ExtractBoodIDFromURL(url)
			if err != nil {
			}
			s.Client.HSet("fidibo", bookID, url)
		} else {

			//e.Request.Visit(e.Request.AbsoluteURL(link))
			err := e.Request.Visit(e.Request.AbsoluteURL(link))
			if err != nil {
				if strings.Contains(link, "page=") {
					fmt.Println(err, e.Request.AbsoluteURL(link))
				}
			}
		}
	})

	err := c.Visit("https://fidibo.com/")
	if err != nil {
		fmt.Println(err)
	}

	c.OnRequest(func(request *colly.Request) {
		//fmt.Println("visiting ", request.URL)
	})
	c.Wait()

}
