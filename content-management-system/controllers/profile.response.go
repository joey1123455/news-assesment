package controllers

import "github.com/joey1123455/news-aggregator-service/content-management-system/models"

type EgFilteredRes struct {
	Username   string           `json:"user_name"`
	Prefrences models.Prefrence `json:"prefrence"`
}

type UpdateProfile struct {
	Status  string        `json:"status"`
	Profile EgFilteredRes `json:"profile"`
}

type CreateProfile struct {
	Status  string        `json:"status"`
	Profile EgFilteredRes `json:"profile"`
}
