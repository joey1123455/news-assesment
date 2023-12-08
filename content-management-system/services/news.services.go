package services

import (
	"context"

	"github.com/joey1123455/news-aggregator-service/content-management-system/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ArticleServices interface {
	Search(key string) ([]models.MongoArticle, error)
	NewsFeed(categories []string, limit, page int) ([]models.MongoArticle, error)
}

type ArticleServiceImp struct {
	ctx        context.Context
	collection *mongo.Collection
}

func NewArticleService(ctx context.Context, collection *mongo.Collection) ArticleServices {
	return &ArticleServiceImp{
		ctx:        ctx,
		collection: collection,
	}
}

func (as *ArticleServiceImp) NewsFeed(categories []string, limit, page int) ([]models.MongoArticle, error) {
	var filter bson.M

	// If categories is not nil, include the category filter
	if categories != nil {
		filter = bson.M{"category": bson.M{"$in": categories}}
	}

	// Define the options to sort by createdAt in descending order, skip, and limit
	options := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(int64((page - 1) * limit)).
		SetLimit(int64(limit))

	// Find articles that match the filter, sort them, skip, and limit
	cursor, err := as.collection.Find(as.ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(as.ctx)

	// Decode the results into a slice of articles
	var articles []models.MongoArticle
	err = cursor.All(as.ctx, &articles)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func (as ArticleServiceImp) Search(key string) ([]models.MongoArticle, error) {
	// Define the filter to search for articles with title or content containing the query
	filter := bson.D{{Key: "$text", Value: bson.D{{Key: "$search", Value: key}}}}

	// Find articles that match the filter
	cursor, err := as.collection.Find(as.ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(as.ctx)

	// Decode the results into a slice of articles
	var articles []models.MongoArticle
	err = cursor.All(as.ctx, &articles)
	if err != nil {
		return nil, err
	}

	return articles, nil
}
