package services

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
	"github.com/joey1123455/news-aggregator-service/news-ags/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ArticleSaverService interface {
	retrieveAllArticles() ([]models.Article, error)
	compareArticles() error
	SaveArticles() error
}

type ArticleSaverServiceImp struct {
	ctx               context.Context
	rClient           *redis.Client
	articleCollection *mongo.Collection
	articleLists      []any
}

func NewArticleSaver(cont context.Context, redDB *redis.Client, monDB *mongo.Collection) ArticleSaverService {
	return &ArticleSaverServiceImp{
		ctx:               cont,
		rClient:           redDB,
		articleCollection: monDB,
		articleLists:      make([]any, 0),
	}
}

// RetrieveAllArticles retrieves all articles from the Redis cache
func (aSS ArticleSaverServiceImp) retrieveAllArticles() ([]models.Article, error) {
	// Key pattern for articles in Redis
	keyPattern := "article:*"

	// Get all keys matching the pattern
	keys, err := aSS.rClient.Keys(context.Background(), keyPattern).Result()
	if err != nil {
		return nil, err
	}

	// Initialize a slice to store articles
	var articles []models.Article

	// Iterate over keys and get values
	for _, key := range keys {
		// Get the article JSON from Redis
		jsonStr, err := aSS.rClient.Get(context.Background(), key).Result()
		if err != nil {
			return nil, err
		}

		// Unmarshal JSON into Article struct
		var article models.Article
		err = json.Unmarshal([]byte(jsonStr), &article)
		if err != nil {
			return nil, err
		}

		// Append the article to the slice
		articles = append(articles, article)
	}

	return articles, nil
}

func (aSS ArticleSaverServiceImp) compareArticles() error {
	return nil
}

func (aSS ArticleSaverServiceImp) SaveArticles() error {
	// articles := interface
	aSS.articleCollection.InsertMany(aSS.ctx, aSS.articleLists)

	// Assume uc.collection is a *mongo.Collection

	// Index model for categories
	categoriesIndex := mongo.IndexModel{
		Keys:    bson.M{"category": 1},            // Assuming "categories" is the field in your documents
		Options: options.Index().SetUnique(false), // SetUnique(false) if categories are not unique
	}

	// Index model for keywords
	keywordsIndex := mongo.IndexModel{
		Keys:    bson.M{"keywords": 1},            // Assuming "keywords" is the field in your documents
		Options: options.Index().SetUnique(false), // SetUnique(false) if keywords are not unique
	}

	// Create indexes
	if _, err := aSS.articleCollection.Indexes().CreateMany(aSS.ctx, []mongo.IndexModel{categoriesIndex, keywordsIndex}); err != nil {
		return err
	}

	return nil
}
