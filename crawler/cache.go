package crawler

import (
	"github.com/matinkhosravani/fidibo_crawler/model"
	"time"
)

type CrawlerCache interface {
	SetCategories(cs []model.Category, expiration time.Duration) error
	GetCategories() (cs []model.Category)
	SetBooksOfCategoryPage(category string, page int, bs []model.Book, expiration time.Duration) error
	BooksOfCategoryPageExists(category string, page int) bool
}
