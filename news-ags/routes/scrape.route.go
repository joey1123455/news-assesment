package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/joey1123455/news-aggregator-service/news-ags/controllers"
	"github.com/joey1123455/news-aggregator-service/news-ags/services"
)

type ScrapeRouteController struct {
	scraperController controllers.ArticleScrapperController
}

func NewScrapeRouteController(ac controllers.ArticleScrapperController) ScrapeRouteController {
	return ScrapeRouteController{
		scraperController: ac,
	}
}

func (rc ScrapeRouteController) ScrapeRoute(rg *gin.RouterGroup, service services.ScrapeArticleService) {
	router := rg.Group("/scrape")

	router.GET("/news", rc.scraperController.ScrapeNews)
}
