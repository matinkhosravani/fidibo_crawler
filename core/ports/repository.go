package ports

import (
	"github.com/matinkhosravani/fidibo_crawler/core/domain"
)

type CrawlerRepository interface {
	Store(b []domain.Book) error
	GetByID(ID string) (domain.Book, bool)
	AddAuthors(bookID string, author []domain.Author)
	AddPublishers(bookID string, publisher []domain.Publisher)
	AddTranslators(bookID string, translators []domain.Translator)
	AddNarrators(bookID string, narrators []domain.Narrator)
}
