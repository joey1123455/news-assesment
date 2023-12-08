package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joey1123455/news-aggregator-service/news-ags/services"
)

type ArticleSaverController struct {
	saverService services.ArticleSaverService
}

func NewArticleSaverController(aSS services.ArticleSaverService) ArticleSaverController {
	return ArticleSaverController{saverService: aSS}
}

func (aSC ArticleSaverController) SaveArticles(ctx *gin.Context) {
	err := aSC.saverService.SaveArticles()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "News scrapped and added to cache"})

}
