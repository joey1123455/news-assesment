package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/joey1123455/news-aggregator-service/news-ags/controllers"
	"github.com/joey1123455/news-aggregator-service/news-ags/services"
)

type SaveRouteController struct {
	saverController controllers.ArticleSaverController
}

func NewSaverRouteController(ac controllers.ArticleSaverController) SaveRouteController {
	return SaveRouteController{
		saverController: ac,
	}
}

func (rc SaveRouteController) SaveRoute(rg *gin.RouterGroup, service services.ArticleSaverService) {
	router := rg.Group("/save")

	router.GET("/news", rc.saverController.SaveArticles)
}
