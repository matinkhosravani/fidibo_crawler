package ports

import (
	"github.com/matinkhosravani/fidibo_crawler/core/domain"
)

type CrawlerRepository interface {
	Store(b []domain.Book) error
}
