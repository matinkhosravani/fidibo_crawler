package main

import (
	"crypto/tls"
	"fmt"
	"github.com/gocolly/colly/v2"
	"net/http"
	"regexp"
)

func main() {
	// Instantiate default collector
	c := colly.NewCollector(
		colly.URLFilters(
			regexp.MustCompile(".*((fidibo\\.com/?$)|(fidibo\\.com/book.*$))"),
		),
	)
	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	})
	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Visit link found on page
		// Only those links are visited which are matched by  any of the URLFilter regexps
		c.Visit(e.Request.AbsoluteURL(link))
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
	// Start scraping on https://fidibo.com
	err := c.Visit("https://fidibo.com/")
	if err != nil {
		fmt.Println(err)
	}
}
