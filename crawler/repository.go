package crawler

import "github.com/matinkhosravani/fidibo_crawler/model"

type CrawlerRepository interface {
	Store(b []model.Book) error
}
