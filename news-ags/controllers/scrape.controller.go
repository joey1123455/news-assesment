package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joey1123455/news-aggregator-service/news-ags/services"
)

type ArticleScrapperController struct {
	scraperService services.ScrapeArticleService
}

func NewArticleScrapperController(sc services.ScrapeArticleService) ArticleScrapperController {
	return ArticleScrapperController{scraperService: sc}
}

// @Summary Scrape News
// @Description Makes a get request to newsapi.io and cahes news objects in database
// @Produce json
// @Success 200 {object} string "News scrapped and added to cache"
// @Failure 500 {object} string "error message"
// @Router /scrape/news [get]
func (aSC ArticleScrapperController) ScrapeNews(ctx *gin.Context) {
	err := aSC.scraperService.ParseArticle()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "News scrapped and added to cache"})

}
