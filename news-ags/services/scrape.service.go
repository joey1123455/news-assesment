package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/go-resty/resty/v2"
	"github.com/joey1123455/news-aggregator-service/news-ags/models"
)

var newsResponsePool = sync.Pool{
	New: func() interface{} {
		return &models.NewsResponse{}
	},
}

type ScrapeArticleService interface {
	ParseArticle() error
	getNews(apiKey, category, nextPage string) error
	cacheArticle(article *models.Article, key string) error
}

type ScrapeArticleServiceImp struct {
	ctx     context.Context
	rClient *redis.Client
	apikey  string
}

func NewScrapper(ctx context.Context, client *redis.Client, key string) ScrapeArticleService {
	return &ScrapeArticleServiceImp{
		ctx:     ctx,
		rClient: client,
		apikey:  key,
	}
}

func (as ScrapeArticleServiceImp) ParseArticle() error {
	var wg sync.WaitGroup
	categories := []string{"business", "entertainment", "health", "science", "sports", "technology", "politics", "tourism", "environment", "domestic"}
	ch := make(chan error, len(categories)+1)
	for _, cat := range categories {
		wg.Add(1)
		go func(category string) {
			defer wg.Done()
			err := as.getNews(as.apikey, category, "")
			if err != nil {
				ch <- err
			}

		}(cat)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()

	// Check for errors from goroutines
	for err := range ch {
		if err != nil {
			// Return the first error received
			return err
		}
	}

	return nil // No errors were received

}

func (as ScrapeArticleServiceImp) getNews(apiKey, category, nextPage string) error {
	result := newsResponsePool.Get().(*models.NewsResponse)
	defer newsResponsePool.Put(result)
	client := resty.New()
	url := "https://newsdata.io/api/1/news"

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	resp, err := client.R().
		// SetContext(ctx).
		SetQueryParam("language", "en").
		SetQueryParam("apiKey", apiKey).
		SetQueryParam("category", category).
		SetQueryParam("full_content", "1").
		SetQueryParam("prioritydomain", "top").
		SetQueryParam("timeframe", "30m").
		Get(url)

	if err != nil {
		return err
	}

	switch resp.StatusCode() {
	case 429:
		return errors.New("api rate limit reached")
	case 500:
		return errors.New("external api failure")
	case 415, 422:
		return errors.New("malformed api querry")
	case 409:
		return errors.New("duplicate api querry")
	case 403:
		return errors.New("corse api error")
	case 400:
		return errors.New("api param missing")
	case 401:
		return errors.New("incorrect api key")
	}

	// Unmarshal the JSON response into the struct
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return err
	}

	for _, article := range result.Articles {
		as.cacheArticle(&article, "articleKey:"+article.Id.Hex())
	}

	return nil
}

func (as ScrapeArticleServiceImp) cacheArticle(article *models.Article, key string) error {
	articleJSON, err := json.Marshal(article)
	if err != nil {
		return fmt.Errorf("error marshaling article: %w", err)
	}

	err = as.rClient.SetNX(key, articleJSON, 2*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("error storing article in Redis: %w", err)
	}

	return nil
}
