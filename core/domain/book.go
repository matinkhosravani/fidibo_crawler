package domain

type Book struct {
	ID            string   `json:"id"`
	Source        string   `json:"-"`
	Title         string   `json:"title"`
	SubTitle      string   `json:"sub_title"`
	Slug          string   `json:"slug"`
	PublishDate   string   `json:"publish_date"`
	Language      string   `json:"language"`
	Free          string   `json:"free"`
	Price         string   `json:"price"`
	Description   string   `json:"description"`
	Format        string   `json:"format"`
	Subscriptions bool     `json:"subscriptions"`
	URL           string   `json:"url"`
	ImageURL      string   `json:"image_url"`
	AudioFormat   bool     `json:"audio_format"`
	Authors       []Author `json:"-"`
}
