package services

import (
	"context"
	"encoding/json"
	"strings"
	"sync"

	"github.com/go-redis/redis"
	"github.com/joey1123455/news-aggregator-service/news-ags/models"
	"github.com/joey1123455/news-aggregator-service/news-ags/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ArticleSaverService interface {
	retrieveAllArticles() ([]models.Article, error)
	compareArticles([]models.Article) ([]models.Article, error)
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
	keys, err := aSS.rClient.Keys(keyPattern).Result()
	if err != nil {
		return nil, err
	}

	// Initialize a slice to store articles
	var articles []models.Article

	// Iterate over keys and get values
	for _, key := range keys {
		// Get the article JSON from Redis
		jsonStr, err := aSS.rClient.Get(key).Result()
		if err != nil {
			return nil, err
		}

		var article models.Article
		err = json.Unmarshal([]byte(jsonStr), &article)
		if err != nil {
			return nil, err
		}

		articles = append(articles, article)
	}

	return articles, nil
}

func (aSS ArticleSaverServiceImp) compareArticles(lst []models.Article) ([]models.Article, error) {
	result := make([]models.Article, len(lst))
	copy(result, lst)

	for i := 0; i < len(lst); i++ {
		for j := i + 1; j < len(lst); j++ {
			if i == j {
				continue
			}
			if strings.Join(lst[i].Keywords, ":") != strings.Join(lst[j].Keywords, ":") {
				continue
			}

			score := utils.CalculateSimilarity(lst[i].Content, lst[j].Content)
			// Check if the similarity score is greater than 0.8
			if score > 0.008988012168019019 {
				// Remove the article with the higher weight
				if result[i].Weight > result[j].Weight {
					result = append(result[:j], result[j+1:]...)
				} else {
					result = append(result[:i], result[i+1:]...)
				}
			}

		}
	}

	return nil, nil
}

func (aSS ArticleSaverServiceImp) SaveArticles() error {
	var wg sync.WaitGroup
	categoryMap := make(map[string][]models.Article)
	articles, err := aSS.retrieveAllArticles()
	if err != nil {
		return err
	}

	for _, article := range articles {
		key := strings.Join(article.Category, " : ")
		categoryMap[key] = append(categoryMap[key], article)
	}
	ch := make(chan []models.Article, len(categoryMap)+1)

	for _, catArticles := range categoryMap {
		wg.Add(1)
		go func(l []models.Article) {
			defer wg.Done()
			res, err := aSS.compareArticles(l)
			if err != nil {
				ch <- res
			}

		}(catArticles)
	}

	for articles := range ch {
		aSS.articleLists = append(aSS.articleLists, articles)
	}

	_, err = aSS.articleCollection.InsertMany(aSS.ctx, aSS.articleLists)
	if err != nil {
		return err
	}

	// Index model for categories
	categoriesIndex := mongo.IndexModel{
		Keys:    bson.M{"category": 1},
		Options: options.Index().SetUnique(false),
	}

	// Index model for keywords
	contentIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "content", Value: "text"}},
	}

	// Index for title search
	// Index model for categories
	titleIndex := mongo.IndexModel{Keys: bson.D{{Key: "title", Value: "text"}}}

	// Create indexes
	if _, err := aSS.articleCollection.Indexes().CreateMany(aSS.ctx, []mongo.IndexModel{categoriesIndex, contentIndex, titleIndex}); err != nil {
		return err
	}

	return nil
}
