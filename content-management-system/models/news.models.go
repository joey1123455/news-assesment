package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Article struct {
	Id          primitive.ObjectID `json:"article_id" bson:"_id" binding:"required"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	URL         string             `json:"link" bson:"link"`
	Source      string             `json:"source_id" bson:"source_id"`
	Weight      int                `json:"source_priority" bson:"source_priority"`
	Keywords    []string           `json:"keywords" bson:"keywords"`
	Author      []string           `json:"creator" bson:"creator"`
	Image       string             `json:"image_url" bson:"image_url"`
	Content     string             `json:"content" bson:"content"`
	Country     []string           `json:"country" bson:"country"`
	Category    []string           `json:"category" bson:"category"`
	Date        string             `json:"pubDate" bson:"pubDate"`
}

type MongoArticle struct {
	ID        primitive.ObjectID `json:"id" bson:"_id" binding:"required"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at" binding:"required"`
	Article   Article            `json:"article" bson:"article" binding:"required"`
}
