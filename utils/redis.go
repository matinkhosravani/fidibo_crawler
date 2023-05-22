package utils

import (
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/redisstorage"
	"log"
)

func SetupCollyStorage(c *colly.Collector) *redisstorage.Storage {
	s := &redisstorage.Storage{
		Address:  "127.0.0.1:6379",
		Password: "",
		DB:       0,
		Prefix:   "fidibo",
	}
	// add repository to the collector
	err := c.SetStorage(s)
	if err != nil {
		panic(err)
	}
	//
	//// delete previous data from repository
	if err := s.Clear(); err != nil {
		log.Fatal(err)
	}

	return s
}
