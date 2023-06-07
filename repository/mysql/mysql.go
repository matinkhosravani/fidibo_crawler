package mysql

import (
	"fmt"
	"github.com/matinkhosravani/fidibo_crawler/core/domain"
	"github.com/matinkhosravani/fidibo_crawler/core/ports"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

type Repository struct {
	DB *gorm.DB
}

func (m Repository) GetByID(ID string) (domain.Book, bool) {
	var b Book
	res := m.DB.Where("book_id = ?", ID).First(&b)
	if res.RowsAffected <= 0 {
		return toDomain(b), false
	}
	return toDomain(b), true
}

func (m Repository) AddAuthors(bookID string, as []domain.Author) {
	var b Book
	m.DB.Where("book_id = ?", bookID).First(&b)

	var gas []Author
	for _, a := range as {
		var ga Author
		var gp Person

		m.DB.FirstOrCreate(&gp, Person{
			SourceID:   a.ID,
			SourceName: "fidibo",
			Name:       a.Name,
		})

		m.DB.FirstOrCreate(&ga, Author{
			Person:   gp,
			PersonID: gp.ID,
		})

		gas = append(gas, ga)
	}
	m.DB.Omit("Authors.*").Model(&b).Association("Authors").Append(gas)
}

func (m Repository) AddPublishers(bookID string, publishers []domain.Publisher) {
	var b Book
	m.DB.Where("book_id = ?", bookID).First(&b)

	var gts []Publisher
	for _, n := range publishers {
		var gt Publisher
		m.DB.FirstOrCreate(&gt, Publisher{
			PublisherID: n.ID,
			Name:        n.Name,
		})

		gts = append(gts, gt)
	}
	m.DB.Omit("Publishers.*").Model(&b).Association("Publishers").Append(gts)
}

func (m Repository) AddTranslators(bookID string, translators []domain.Translator) {
	var b Book
	m.DB.Where("book_id = ?", bookID).First(&b)

	var gts []Translator
	for _, n := range translators {
		var gt Translator
		var gp Person

		m.DB.FirstOrCreate(&gp, Person{
			SourceID:   n.ID,
			SourceName: "fidibo",
			Name:       n.Name,
		})

		m.DB.FirstOrCreate(&gt, Translator{
			Person:   gp,
			PersonID: gp.ID,
		})

		gts = append(gts, gt)
	}
	m.DB.Omit("Translators.*").Model(&b).Association("Translators").Append(gts)
}
func (m Repository) AddNarrators(bookID string, narrators []domain.Narrator) {
	var b Book
	m.DB.Where("book_id = ?", bookID).First(&b)

	var gns []Narrator
	for _, n := range narrators {
		var gn Narrator
		var gp Person

		m.DB.FirstOrCreate(&gp, Person{
			SourceID:   n.ID,
			SourceName: "fidibo",
			Name:       n.Name,
		})

		m.DB.FirstOrCreate(&gn, Translator{
			Person:   gp,
			PersonID: gp.ID,
		})

		gns = append(gns, gn)
	}
	m.DB.Omit("Narrators.*").Model(&b).Association("Narrators").Append(gns)
}

func (m Repository) Store(bs []domain.Book) error {
	gormBooks := fromDomainSlice(bs)

	for _, gb := range gormBooks {
		m.DB.Create(&gb)
	}

	return nil
}

func fromDomainSlice(bs []domain.Book) []Book {
	var gbs []Book
	for _, b := range bs {
		gb := fromDomain(b)
		gbs = append(gbs, gb)
	}

	return gbs
}

func fromDomain(b domain.Book) Book {
	gb := Book{
		BookID:      b.ID,
		Title:       b.Title,
		SubTitle:    b.SubTitle,
		Slug:        b.Slug,
		PublishDate: b.PublishDate,
		Language:    b.Language,
		Free:        b.Free,
		Price:       b.Price,
		Description: b.Description,
		Format:      b.Format,
		Source:      b.Source,
		URL:         b.URL,
		ImageURL:    b.ImageURL,
		AudioFormat: b.AudioFormat,
	}
	var gAuthors []Author

	for _, author := range b.Authors {
		gAuthor := Author{
			Person: Person{
				SourceID:   author.ID,
				SourceName: "fidibo",
				Name:       author.Name,
			},
		}
		gAuthors = append(gAuthors, gAuthor)
	}
	gb.Authors = gAuthors
	return gb
}

func toDomainSlice(bs []Book) []domain.Book {
	var gbs []domain.Book
	for _, b := range bs {
		gb := toDomain(b)
		gbs = append(gbs, gb)
	}

	return gbs
}

func toDomain(b Book) domain.Book {
	gb := domain.Book{
		ID:          b.BookID,
		Title:       b.Title,
		SubTitle:    b.SubTitle,
		Slug:        b.Slug,
		PublishDate: b.PublishDate,
		Language:    b.Language,
		Free:        b.Free,
		Price:       b.Price,
		Source:      b.Source,
		Description: b.Description,
		Format:      b.Format,
		URL:         b.URL,
		ImageURL:    b.ImageURL,
		AudioFormat: b.AudioFormat,
	}
	var gAuthors []domain.Author

	for _, author := range b.Authors {
		gAuthor := domain.Author{
			ID:   author.Person.SourceID,
			Name: author.Person.Name,
		}
		gAuthors = append(gAuthors, gAuthor)
	}
	gb.Authors = gAuthors
	return gb
}
func NewRepository() (ports.CrawlerRepository, error) {
	repo := &Repository{}

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"))

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&Book{}, &Person{}, &Author{})
	if err != nil {
		return nil, err
	}

	repo.DB = db

	return repo, nil
}
