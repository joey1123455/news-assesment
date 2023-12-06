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

func (aSC ArticleScrapperController) ScrapeNews(ctx *gin.Context) {
	err := aSC.scraperService.ParseArticle()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "News scrapped and added to cache"})

}
