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

// @Summary Save News
// @Description Saves news articles stored in a redis cache into a mongo collection
// @Produce json
// @Success 200 {object} string "News scrapped and added to cache"
// @Failure 500 {object} string "error message"
// @Router /save/news [get]
func (aSC ArticleSaverController) SaveArticles(ctx *gin.Context) {
	err := aSC.saverService.SaveArticles()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "News scrapped and added to cache"})

}
