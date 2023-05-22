package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/matinkhosravani/fidibo_crawler/app"
	"github.com/matinkhosravani/fidibo_crawler/cache/redis"
	"github.com/matinkhosravani/fidibo_crawler/crawler"
	"github.com/matinkhosravani/fidibo_crawler/model"
	"github.com/matinkhosravani/fidibo_crawler/repository/mongo"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	concurrency = 7
	semaChan    = make(chan struct{}, concurrency)
)
var c = crawler.NewCrawler()

func main() {
	app.LoadEnv()
	c.Repo = chooseRepo()
	c.Cache = chooseCache()

	log.Println("Connected to MongoDB!")
	categories := findCategories()
	log.Println("Founded all Root categories")

	done := make(chan struct{})

	for _, category := range categories {
		responses := make(chan CategoryResponse)
		totalPages, _ := findTotalPages(category)
		go func() {
			for i := 0; i < totalPages; i++ {
				resp := <-responses
				err := c.Repo.Store(resp.Books)
				if err != nil {
					log.Fatal(err.Error())
				}
			}
		}()

		for i := 1; i <= totalPages; i++ {
			isPageCached := c.Cache.BooksOfCategoryPageExists(category.Name, i)
			if !isPageCached {
				semaChan <- struct{}{} // block while full
				go getBooksByCategory(category, i, responses)
			}
		}
	}

	go func() {
		defer close(done)
	}()

	<-done
}

func findCategories() []model.Category {
	categories := c.Cache.GetCategories()
	if categories != nil {
		return categories
	}
	return getAllRootCategories()
}

func findTotalPages(category model.Category) (int, error) {
	resChannel := make(chan CategoryResponse)
	go getBooksByCategory(category, 1, resChannel)
	resp := <-resChannel

	return resp.TotalPages, nil
}

type CategoryResponse struct {
	Books       []model.Book  `json:"books"`
	Page        int           `json:"page"`
	PerPage     int           `json:"size"`
	Sorting     string        `json:"sorting"`
	Keyword     string        `json:"keyword"`
	BookFormats []interface{} `json:"book_formats"`
	Total       int           `json:"total"`
	TotalPages  int           `json:"total_pages"`
}

type option func(url *string)

func getBooksByCategory(category model.Category, page int, responseStream chan<- CategoryResponse, options ...option) {
	defer func() {
		<-semaChan // read releases a slot
	}()
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
	var categoryResp CategoryResponse
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

func getAllRootCategories() []model.Category {
	var categories []model.Category

	c := colly.NewCollector(
		colly.AllowedDomains("fidibo.com"),
	)
	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	})
	// On every an element which has href attribute call callback
	c.OnHTML("ul.dropdown-menu > div > li > a", func(e *colly.HTMLElement) {
		category := strings.Replace(e.Attr("href"), "/category/", "", 1)
		if category != "" {
			categories = append(categories, model.Category{
				Name: category,
			})
		}
	})

	err := c.Visit("https://fidibo.com/")
	if err != nil {
		fmt.Println(err)
	}

	return categories
}

func chooseCache() crawler.CrawlerCache {
	switch os.Getenv("CACHE_NAME") {
	case "redis":
		cache, err := redis.NewRedisCacheRepository()
		if err != nil {
			log.Fatal(err.Error())
		}
		return cache
	}
	return nil
}

func chooseRepo() crawler.CrawlerRepository {
	switch os.Getenv("DB_NAME") {
	case "mongo":
		repo, err := mongo.NewMongoRepository()
		if err != nil {
			log.Fatal(err.Error())
		}
		return repo
	}

	return nil
}
