package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/matinkhosravani/fidibo_crawler/app"
	"github.com/matinkhosravani/fidibo_crawler/storage"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
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

func main() {
	app.LoadEnv()
	storage.ConnectToMongo()
	log.Println("Connected to MongoDB!")
	categories := findCategories()
	log.Println("Founded all Root categories")

	collection := storage.Client.Database(os.Getenv("MONGO_DATABASE")).Collection("books")
	done := make(chan struct{})

	for _, category := range categories {
		responses := make(chan CategoryResponse)
		totalPages, _ := findTotalPages(category)
		go persistBooks(responses, collection, totalPages)

		for i := 1; i <= totalPages; i++ {
			page := storage.Redis().Get(context.Background(), fmt.Sprintf("%s_%d", category, i))
			if page.Err() == redis.Nil {
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

func persistBooks(responses <-chan CategoryResponse, collection *mongo.Collection, totalPages int) {
	for i := 0; i < totalPages; i++ {
		resp := <-responses
		for _, book := range resp.Books {
			_, err := collection.InsertOne(context.Background(), book)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func findCategories() []string {
	cachedCategories := storage.Redis().Get(context.Background(), "fidibo_categories")
	var categories []string
	if cachedCategories.Err() != redis.Nil {
		json.Unmarshal([]byte(cachedCategories.Val()), &categories)
	} else {
		categories = getAllRootCategories()
		j, err := json.Marshal(categories)
		if err != nil {
			fmt.Println(err)
		}
		s := storage.Redis().Set(context.Background(), "fidibo_categories", string(j), time.Hour*1)
		if s.Err() != redis.Nil {
			fmt.Println(s.Err().Error())
		}
	}
	return categories
}

func findTotalPages(category string) (int, error) {
	resChannel := make(chan CategoryResponse)
	go getBooksByCategory(category, 1, resChannel)
	resp := <-resChannel

	return resp.TotalPages, nil
}

type Book struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	SubTitle      string   `json:"sub_title"`
	Slug          string   `json:"slug"`
	PublishDate   string   `json:"publish_date"`
	Language      string   `json:"language"`
	Free          string   `json:"free"`
	Price         string   `json:"price"`
	Description   string   `json:"description"`
	PublisherID   string   `json:"publisher_id"`
	TranslatorID  string   `json:"translator_id"`
	NarratorID    any      `json:"narrator_id"`
	Format        string   `json:"format"`
	Subscriptions bool     `json:"subscriptions"`
	URL           string   `json:"url"`
	ImageURL      string   `json:"image_url"`
	AudioFormat   bool     `json:"audio_format"`
	Authors       []Author `json:"authors"`
}

type Author struct {
	ID   string
	Name string
}

type Narrotor struct {
	ID   string
	Name string
}

type Publisher struct {
	ID string
}

type CategoryResponse struct {
	Books       []Book        `json:"books"`
	Page        int           `json:"page"`
	PerPage     int           `json:"size"`
	Sorting     string        `json:"sorting"`
	Keyword     string        `json:"keyword"`
	BookFormats []interface{} `json:"book_formats"`
	Total       int           `json:"total"`
	TotalPages  int           `json:"total_pages"`
}

type option func(url *string)

func getBooksByCategory(category string, page int, responseStream chan<- CategoryResponse, options ...option) {
	defer func() {
		<-semaChan // read releases a slot
	}()
	url := fmt.Sprintf("https://fidibo.com/category/%v?page=%v", category, page)
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
	j, err := json.Marshal(categoryResp.Books)
	if err != nil {
		fmt.Println(err)
	}
	storage.Redis().Set(context.Background(), fmt.Sprintf("%s_%d", category, page), string(j), time.Hour*24)
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

func getAllRootCategories() []string {
	var categories []string

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
			categories = append(categories, category)
		}
	})

	err := c.Visit("https://fidibo.com/")
	if err != nil {
		fmt.Println(err)
	}

	return categories
}
