package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/matinkhosravani/fidibo_crawler/crawler"
	"github.com/matinkhosravani/fidibo_crawler/model"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"os"
	"time"
)

type redisCacheRepository struct {
	client *redis.Client
}

func (r redisCacheRepository) SetCategories(cs []model.Category, expiration time.Duration) error {
	j, err := json.Marshal(cs)
	if err != nil {
		fmt.Println(err)
	}
	s := r.client.Set(context.Background(), generateCategoriesKey(), string(j), expiration)
	if s.Err() != redis.Nil {
		return errors.Wrap(err, "repository.SetCategories")
	}

	return nil
}

func generateCategoriesKey() string {
	return "fidibo_categories"
}

func (r redisCacheRepository) GetCategories() []model.Category {
	cache := r.client.Get(context.Background(), generateCategoriesKey())
	var categories []model.Category
	if cache.Err() != redis.Nil {
		json.Unmarshal([]byte(cache.Val()), &categories)
	}

	return categories
}

func (r redisCacheRepository) SetBooksOfCategoryPage(category string, page int, bs []model.Book, expiration time.Duration) error {
	j, err := json.Marshal(bs)
	if err != nil {
		return err
	}
	r.client.Set(context.Background(), generateCategoryPageKey(category, page), string(j), expiration)

	return nil
}

func generateCategoryPageKey(category string, page int) string {
	return fmt.Sprintf("%s_%d", category, page)
}

func (r redisCacheRepository) BooksOfCategoryPageExists(category string, page int) bool {
	res := r.client.Get(context.Background(), fmt.Sprintf("%s_%v", category, page))
	if res.Err() == redis.Nil {
		return false
	}

	return true

}

func NewRedisCacheRepository() (crawler.CrawlerCache, error) {
	repo := &redisCacheRepository{}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	repo.client = client

	return repo, nil

}
func Redis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
