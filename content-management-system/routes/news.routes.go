package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/joey1123455/news-aggregator-service/content-management-system/controllers"
	"github.com/joey1123455/news-aggregator-service/content-management-system/services"
)

type NewsRouteController struct {
	newsController controllers.NewsController
}

func NewNewsControllerRoute(newsController controllers.NewsController) NewsRouteController {
	return NewsRouteController{newsController}
}

func (r *NewsRouteController) NewsRoute(rg *gin.RouterGroup, service services.ArticleServices) {
	router := rg.Group("/news")

	router.GET("/feed", r.newsController.Feed)
	router.GET("/search", r.newsController.Search)
}
