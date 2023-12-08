package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joey1123455/news-aggregator-service/content-management-system/services"
)

type NewsController struct {
	service services.ArticleServices
}

func NewNewsController(service services.ArticleServices) NewsController {
	return NewsController{
		service: service,
	}
}

// @Summary Feed
// @Description Returns a news feed according to users prefrences.
// @Security ApiKeyAuth
// @Produce json
// @Param page query string true "amount of page results to return"
// @Param limit query string true "limit per page"
// @Success 201 {object} SearchResponse
// @Failure 404 {object} string "querry not passed"
// @Failure 502 {object} string "error message"
// @Router /news/feed [get]
func (nc NewsController) Feed(ctx *gin.Context) {
	prefrence := ctx.MustGet("currentUserPrefrence").([]string)
	page := ctx.DefaultQuery("page", "1")
	limit := ctx.DefaultQuery("limit", "10")

	intPage, err := strconv.Atoi(page)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	intLimit, err := strconv.Atoi(limit)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	posts, err := nc.service.NewsFeed(prefrence, intLimit, intPage)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(posts), "articles": posts})
}

// @Summary Search
// @Description Searches database for articles containing the required keywords.
// @Security ApiKeyAuth
// @Produce json
// @Param q query string true "Search query"
// @Success 201 {object} SearchResponse
// @Failure 404 {object} string "querry not passed"
// @Failure 502 {object} string "error message"
// @Router /news/search [get]
func (nc NewsController) Search(ctx *gin.Context) {
	query := ctx.Query("q")
	if query == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "querry not passed"})
	}

	articles, err := nc.service.Search(query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "failed", "articles": articles})
}
