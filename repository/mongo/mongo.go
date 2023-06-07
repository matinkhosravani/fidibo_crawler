package mongo

import (
	"context"
	"fmt"
	"github.com/matinkhosravani/fidibo_crawler/core/domain"
	"github.com/matinkhosravani/fidibo_crawler/core/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

type Repository struct {
	Client   *mongo.Client
	DB       string
	MongoURL string
	Timeout  int
}

func (m Repository) AddPublishers(bookID string, publishers []domain.Publisher) {
	_, err := m.Client.Database(m.DB).Collection("books").UpdateOne(context.TODO(), bson.D{
		{
			"id", bookID,
		},
	}, bson.D{
		{"$set", bson.D{
			{"publishers", publishers},
		}},
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (m Repository) AddTranslators(bookID string, translators []domain.Translator) {
	_, err := m.Client.Database(m.DB).Collection("books").UpdateOne(context.TODO(), bson.D{
		{
			"id", bookID,
		},
	}, bson.D{
		{"$set", bson.D{
			{"translators", translators},
		}},
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (m Repository) AddNarrators(bookID string, narrators []domain.Narrator) {
	_, err := m.Client.Database(m.DB).Collection("books").UpdateOne(context.TODO(), bson.D{
		{
			"id", bookID,
		},
	}, bson.D{
		{"$set", bson.D{
			{"narrators", narrators},
		}},
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (m Repository) GetByID(ID string) (domain.Book, bool) {
	var b domain.Book

	err := m.Client.Database(m.DB).Collection("books").FindOne(context.TODO(), bson.D{
		{
			"id", ID,
		},
	}).Decode(&b)
	if err == mongo.ErrNoDocuments {
		return b, false
	}
	if err != nil {
		log.Fatal(err)
	}

	return b, true
}
func (m Repository) AddAuthors(bookID string, authors []domain.Author) {
	_, err := m.Client.Database(m.DB).Collection("books").UpdateOne(context.TODO(), bson.D{
		{
			"id", bookID,
		},
	}, bson.D{
		{"$set", bson.D{
			{"authors", authors},
		}},
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (m Repository) Store(bs []domain.Book) error {
	collection := m.Client.Database(m.DB).Collection("books")
	is := make([]interface{}, len(bs))
	for i := range bs {
		is[i] = bs[i]
	}
	_, err := collection.InsertMany(context.Background(), is)
	if err != nil {
		return err
	}

	return nil
}

var Client *mongo.Client

func NewRepository() (ports.CrawlerRepository, error) {

	repo := &Repository{
		Client:   nil,
		DB:       os.Getenv("MONGO_DATABASE"),
		MongoURL: fmt.Sprintf("mongodb://%v:%v", os.Getenv("MONGO_HOST"), os.Getenv("MONGO_PORT")),
	}
	// Set client options
	clientOptions := options.Client().ApplyURI(repo.MongoURL)
	// Connect to MongoDB
	var err error
	Client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = Client.Ping(context.TODO(), nil)

	if err != nil {
		return nil, err
	}

	repo.Client = Client

	return repo, nil
}
