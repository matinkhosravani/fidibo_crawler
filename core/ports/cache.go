package ports

import (
	domain2 "github.com/matinkhosravani/fidibo_crawler/core/domain"
	"time"
)

type CrawlerCache interface {
	SetCategories(cs []domain2.Category, expiration time.Duration) error
	GetCategories() (cs []domain2.Category)
	SetBooksOfCategoryPage(category string, page int, bs []domain2.Book, expiration time.Duration) error
	BooksOfCategoryPageExists(category string, page int) bool
	GetBookURLS() map[string]string
}
