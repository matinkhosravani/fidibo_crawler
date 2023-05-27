package crawler

import "github.com/matinkhosravani/fidibo_crawler/core/ports"

type Crawler struct {
	Repo  ports.CrawlerRepository
	Cache ports.CrawlerCache
}

func NewCrawler() *Crawler {
	c := &Crawler{}
	return c
}
