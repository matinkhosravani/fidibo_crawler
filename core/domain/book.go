package domain

type Book struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	SubTitle      string   `json:"sub_title"`
	Slug          string   `json:"slug"`
	PublishDate   string   `json:"publish_date"`
	Language      string   `json:"language"`
	Free          string   `json:"free"`
	Price         string   `json:"price"`
	Description   string   `json:"description"`
	PublisherID   string   `json:"publisher_id"`
	TranslatorID  string   `json:"translator_id"`
	NarratorID    string   `json:"narrator_id"`
	Format        string   `json:"format"`
	Subscriptions bool     `json:"subscriptions"`
	URL           string   `json:"url"`
	ImageURL      string   `json:"image_url"`
	AudioFormat   bool     `json:"audio_format"`
	Authors       []Author `json:"authors"`
}
