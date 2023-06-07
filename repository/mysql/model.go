package mysql

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	BookID      string       `json:"id"`
	Title       string       `json:"title"`
	SubTitle    string       `json:"sub_title"`
	Slug        string       `json:"slug"`
	PublishDate string       `json:"publish_date"`
	Language    string       `json:"language"`
	Free        string       `json:"free"`
	Price       string       `json:"price"`
	Description string       `json:"description"`
	Publishers  []Publisher  `gorm:"many2many:book_publishers;"`
	Narrators   []Narrator   `gorm:"many2many:book_narrators;"`
	Translators []Translator `gorm:"many2many:book_translators;"`
	Authors     []Author     `json:"authors" gorm:"many2many:book_authors;"`
	Source      string
	Format      string `json:"format"`
	URL         string `json:"url"`
	ImageURL    string `json:"image_url"`
	AudioFormat bool   `json:"audio_format"`
}

type Author struct {
	gorm.Model
	PersonID uint
	Person   Person `gorm:"foreignKey:PersonID"`
	Books    []Book `gorm:"many2many:book_authors;"`
}

type Publisher struct {
	gorm.Model
	PublisherID string `json:"id"`
	Name        string
	Books       []Book `gorm:"many2many:book_publishers;"`
}
type Narrator struct {
	gorm.Model
	PersonID uint
	Person   Person `gorm:"foreignKey:PersonID"`
	Books    []Book `gorm:"many2many:book_narrators;"`
}
type Translator struct {
	gorm.Model
	PersonID uint
	Person   Person `gorm:"foreignKey:PersonID"`
	Books    []Book `gorm:"many2many:book_translators;"`
}

type Person struct {
	gorm.Model
	SourceID   string
	SourceName string
	Name       string
}
