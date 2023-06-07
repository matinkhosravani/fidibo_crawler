package main

import (
	"crypto/tls"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/matinkhosravani/fidibo_crawler/app"
	"github.com/matinkhosravani/fidibo_crawler/cache/redis"
	"github.com/matinkhosravani/fidibo_crawler/core/domain"
	"github.com/matinkhosravani/fidibo_crawler/repository/mysql"
	"github.com/matinkhosravani/fidibo_crawler/utils"
	"log"
	"net/http"
	"time"
)

func main() {
	app.LoadEnv()
	// Instantiate default collector
	c := colly.NewCollector(
		colly.AllowedDomains("fidibo.com"),
	)
	cacheRepo, err := redis.NewCacheRepository()
	repo, err := mysql.NewRepository()
	if err != nil {
		log.Fatal(err.Error())
	}
	books := cacheRepo.GetBookURLS()
	s := utils.SetupCollyStorage(c)
	// close redis client
	defer s.Client.Close()
	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	})

	c.OnHTML("li.author_title", func(e *colly.HTMLElement) {
		t := e.DOM.Find("span.title-list").Text()
		ID, err := utils.ExtractBoodIDFromURL(e.Request.URL.String())
		if err != nil {
			log.Println(err, e.Request.URL.String())
		}
		if _, ok := repo.GetByID(ID); !ok {
			return
		}

		switch t {
		case "نویسنده":
			var as []domain.Author
			e.DOM.Find("a span").Each(func(i int, selection *goquery.Selection) {
				name := selection.Text()
				if id, ok := selection.Parent().Attr("data-ut-object-id"); ok {
					as = append(as, domain.Author{
						ID:   id,
						Name: name,
					})
				}
			})
			repo.AddAuthors(ID, as)
		case "مترجم":
			var ts []domain.Translator
			e.DOM.Find("a span").Each(func(i int, selection *goquery.Selection) {
				name := selection.Text()
				if id, ok := selection.Parent().Attr("data-ut-object-id"); ok {
					ts = append(ts, domain.Translator{
						ID:   id,
						Name: name,
					})
				}
			})
			repo.AddTranslators(ID, ts)
		case "گوینده":
			var as []domain.Narrator
			e.DOM.Find("a span").Each(func(i int, selection *goquery.Selection) {
				name := selection.Text()
				if id, ok := selection.Parent().Attr("data-ut-object-id"); ok {
					as = append(as, domain.Narrator{
						ID:   id,
						Name: name,
					})
				}
			})
			repo.AddNarrators(ID, as)
		}
	})
	c.OnHTML("div.book-tags ul li a", func(e *colly.HTMLElement) {
		ID, err := utils.ExtractBoodIDFromURL(e.Request.URL.String())
		if err != nil {
			log.Fatal(err)
		}
		if _, ok := repo.GetByID(ID); !ok {
			return
		}

		ps := []domain.Publisher{
			{
				ID:   e.Attr("data-ut-object-id"),
				Name: e.Text,
			},
		}

		repo.AddPublishers(ID, ps)
	})
	// On every an element which has href attribute call callback
	for k, v := range books {
		err = c.Visit(v)
		time.Sleep(time.Millisecond * 200)
		if err != nil {
			fmt.Println(err, k, v)
		}
	}
}
