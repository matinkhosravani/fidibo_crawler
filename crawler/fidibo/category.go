package fidibo

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/matinkhosravani/fidibo_crawler/core/domain"
	"github.com/matinkhosravani/fidibo_crawler/crawler"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	DOMAIN = "fidibo.com"
)

var wg sync.WaitGroup
var (
	concurrency = 7
	semaChan    = make(chan struct{}, concurrency)
)

func Crawl(c *crawler.Crawler, booksStream chan<- []domain.Book) {
	categories := findCategories(c)
	for _, category := range categories {
		responses := make(chan categoryResponse)
		totalPages, _ := findTotalPages(c, category)
		wg.Wait()
		go func() {
			for resp := range responses {
				books := resp.Books
				for _, book := range books {
					book.Source = "fidibo"
				}

				booksStream <- books
			}
		}()
		for i := 1; i <= totalPages; i++ {
			isPageCached := c.Cache.BooksOfCategoryPageExists(category.Name, i)
			if isPageCached {
				continue
			}
			semaChan <- struct{}{} // block while full
			wg.Add(1)
			go getBooksByCategory(c, category, i, responses)
		}
		wg.Wait()
		close(responses)
	}
}

func findCategories(c *crawler.Crawler) []domain.Category {
	categories := c.Cache.GetCategories()
	if categories != nil {
		c.Cache.SetCategories(categories, time.Hour*24)
		return categories
	}
	return getAllRootCategories()
}

func findTotalPages(c *crawler.Crawler, category domain.Category) (int, error) {
	resChannel := make(chan categoryResponse)
	wg.Add(1)
	go getBooksByCategory(c, category, 1, resChannel)
	resp := <-resChannel

	return resp.TotalPages, nil
}

type categoryResponse struct {
	Books       []domain.Book `json:"books"`
	Page        int           `json:"page"`
	PerPage     int           `json:"size"`
	Sorting     string        `json:"sorting"`
	Keyword     string        `json:"keyword"`
	BookFormats []interface{} `json:"book_formats"`
	Total       int           `json:"total"`
	TotalPages  int           `json:"total_pages"`
}

type option func(url *string)

func getBooksByCategory(c *crawler.Crawler, category domain.Category, page int, responseStream chan<- categoryResponse, options ...option) {
	defer func() {
		<-semaChan // read releases a slot
	}()
	defer wg.Done()

	url := fmt.Sprintf("https://fidibo.com/category/%v?page=%v", category.Name, page)
	for _, opt := range options {
		opt(&url)
	}
	fmt.Println(url)
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, _ := http.NewRequest("GET", url, nil)
	setHeaders(req, url)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	var categoryResp categoryResponse
	err = json.Unmarshal(body, &categoryResp)
	if err != nil {
		log.Fatal(err.Error())
	}
	responseStream <- categoryResp
	err = c.Cache.SetBooksOfCategoryPage(category.Name, page, categoryResp.Books, time.Hour*24)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func setHeaders(req *http.Request, url string) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/112.0")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", url)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
}

func withParams(key, value string) option {
	return func(url *string) {
		*url += fmt.Sprintf("&%v=%v", key, value)
	}
}

func getAllRootCategories() []domain.Category {
	var categories []domain.Category

	c := colly.NewCollector(
		colly.AllowedDomains(DOMAIN),
	)
	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	})
	// On every an element which has href attribute call callback
	c.OnHTML("ul.dropdown-menu > div > li > a", func(e *colly.HTMLElement) {
		category := strings.Replace(e.Attr("href"), "/category/", "", 1)
		if category != "" {
			categories = append(categories, domain.Category{
				Name: category,
			})
		}
	})

	err := c.Visit(fmt.Sprintf("https://%v/", DOMAIN))
	if err != nil {
		fmt.Println(err)
	}

	return categories
}
