package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/joey1123455/news-aggregator-service/news-ags/models"
)

type ArticleRepository interface {
	SaveArticle(key string, article models.Article) error
	GetArticle(key string) (*models.Article, error)
}

type RedisArticleRepository struct {
	client *redis.Client
	ctx    context.Context
}

// SaveArticle saves an article in Redis
func (r *RedisArticleRepository) SaveArticle(key string, article models.Article) error {
	articleJSON, err := json.Marshal(article)
	if err != nil {
		return fmt.Errorf("error marshaling article: %w", err)
	}

	err = r.client.Set(r.ctx, key, articleJSON, 0).Err()
	if err != nil {
		return fmt.Errorf("error storing article in Redis: %w", err)
	}

	return nil
}

// GetArticle retrieves an article from Redis
func (r *RedisArticleRepository) GetArticle(key string) (*models.Article, error) {
	storedJSON, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("error retrieving article from Redis: %w", err)
	}

	var retrievedArticle models.Article
	err = json.Unmarshal([]byte(storedJSON), &retrievedArticle)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling article: %w", err)
	}

	return &retrievedArticle, nil
}

// NewRedisArticleRepository creates a new instance of RedisArticleRepository
func NewRedisArticleRepository(client *redis.Client) *RedisArticleRepository {
	return &RedisArticleRepository{client: client}
}
