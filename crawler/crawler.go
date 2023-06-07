package crawler

import "github.com/matinkhosravani/fidibo_crawler/core/ports"

type Crawler struct {
	Repo  ports.CrawlerRepository
	Cache ports.CrawlerCache
}

func NewCrawler(cache ports.CrawlerCache, repo ports.CrawlerRepository) *Crawler {
	c := &Crawler{
		Repo:  repo,
		Cache: cache,
	}

	return c
}
