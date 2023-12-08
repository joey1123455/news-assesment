package controllers

import "github.com/joey1123455/news-aggregator-service/content-management-system/models"

type EgFilteredRes struct {
	Username   string           `json:"user_name"`
	Prefrences models.Prefrence `json:"prefrence"`
}

type Profile struct {
	Status  string        `json:"status"`
	Profile EgFilteredRes `json:"profile"`
}

type SearchResponse struct {
	Status   string                `json:"status"`
	Articles []models.MongoArticle `json:"articles"`
}

type FeedResponse struct {
	Status   string                `json:"status"`
	Length   int                   `json:"results"`
	Articles []models.MongoArticle `json:"articles"`
}
