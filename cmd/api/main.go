package main

import (
	"github.com/matinkhosravani/fidibo_crawler/app"
	"github.com/matinkhosravani/fidibo_crawler/cache/redis"
	"github.com/matinkhosravani/fidibo_crawler/core/domain"
	"github.com/matinkhosravani/fidibo_crawler/core/ports"
	"github.com/matinkhosravani/fidibo_crawler/crawler"
	"github.com/matinkhosravani/fidibo_crawler/crawler/fidibo"
	"github.com/matinkhosravani/fidibo_crawler/repository/mongo"
	"github.com/matinkhosravani/fidibo_crawler/repository/mysql"
	"log"
	"os"
)

func main() {

	app.LoadEnv()
	repo := chooseRepo()
	cache := chooseCache()
	c := crawler.NewCrawler(cache, repo)

	booksStream := make(chan []domain.Book)

	go func() {
		for books := range booksStream {
			err := c.Repo.Store(books)
			if err != nil {
				log.Fatal(err.Error())
			}
		}
	}()

	fidibo.Crawl(c, booksStream)

	close(booksStream)
}

func chooseCache() ports.CrawlerCache {
	switch os.Getenv("CACHE_NAME") {
	case "redis":
		cache, err := redis.NewCacheRepository()
		if err != nil {
			log.Fatal(err.Error())
		}
		return cache
	}
	return nil
}

func chooseRepo() ports.CrawlerRepository {
	switch os.Getenv("DB_NAME") {
	case "mongo":
		repo, err := mongo.NewRepository()
		if err != nil {
			log.Fatal(err.Error())
		}
		return repo
	case "mysql":
		repo, err := mysql.NewRepository()
		if err != nil {
			log.Fatal(err.Error())
		}
		return repo
	}

	return nil
}
