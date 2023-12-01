package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/joey1123455/news-aggregator-service/users/controllers"
	"github.com/joey1123455/news-aggregator-service/users/middleware"
	"github.com/joey1123455/news-aggregator-service/users/services"
)

type UserRouteController struct {
	userController controllers.UserController
}

func NewRouteUserController(userController controllers.UserController) UserRouteController {
	return UserRouteController{userController}
}

func (uc *UserRouteController) UserRoute(rg *gin.RouterGroup, userService services.UserService) {

	router := rg.Group("users")
	router.Use(middleware.DeserializeUser(userService))
	router.GET("/me", uc.userController.GetMe)
}
