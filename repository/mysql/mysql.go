package mysql

import (
	"fmt"
	"github.com/matinkhosravani/fidibo_crawler/crawler"
	"github.com/matinkhosravani/fidibo_crawler/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

type Repository struct {
	DB *gorm.DB
}

func (m Repository) Store(bs []model.Book) error {
	gormBooks := adaptModels(bs)

	for _, gb := range gormBooks {
		m.DB.Create(&gb)
	}

	return nil
}

func adaptModels(bs []model.Book) []Book {
	var gbs []Book
	for _, b := range bs {
		gb := Book{
			BookID:        b.ID,
			Title:         b.Title,
			SubTitle:      b.SubTitle,
			Slug:          b.Slug,
			PublishDate:   b.PublishDate,
			Language:      b.Language,
			Free:          b.Free,
			Price:         b.Price,
			Description:   b.Description,
			PublisherID:   b.PublisherID,
			TranslatorID:  b.TranslatorID,
			NarratorID:    b.NarratorID,
			Format:        b.Format,
			Subscriptions: b.Subscriptions,
			URL:           b.URL,
			ImageURL:      b.ImageURL,
			AudioFormat:   b.AudioFormat,
		}
		var gAuthors []Author

		for _, author := range b.Authors {
			gAuthor := Author{
				AuthorID: author.ID,
				Name:     author.Name,
			}
			gAuthors = append(gAuthors, gAuthor)
		}
		gb.Authors = gAuthors
		gbs = append(gbs, gb)
	}

	return gbs
}

func NewRepository() (crawler.CrawlerRepository, error) {
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

	err = db.AutoMigrate(&Book{}, &Author{})
	if err != nil {
		return nil, err
	}

	repo.DB = db

	return repo, nil
}
