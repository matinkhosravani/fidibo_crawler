package storage

import (
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/redisstorage"
	"github.com/redis/go-redis/v9"
	"log"
)

func Redis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func SetupCollyStorage(c *colly.Collector) *redisstorage.Storage {
	s := &redisstorage.Storage{
		Address:  "127.0.0.1:6379",
		Password: "",
		DB:       0,
		Prefix:   "fidibo",
	}
	// add storage to the collector
	err := c.SetStorage(s)
	if err != nil {
		panic(err)
	}
	//
	//// delete previous data from storage
	if err := s.Clear(); err != nil {
		log.Fatal(err)
	}

	return s
}
