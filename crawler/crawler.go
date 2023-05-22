package crawler

type Crawler struct {
	Repo  CrawlerRepository
	Cache CrawlerCache
}

func NewCrawler() *Crawler {
	c := &Crawler{}
	return c
}
