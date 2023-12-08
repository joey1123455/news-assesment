package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/joey1123455/news-aggregator-service/content-management-system/controllers"
	"github.com/joey1123455/news-aggregator-service/content-management-system/services"
)

type ProfileRouteController struct {
	profileController controllers.ProfileController
}

func NewprofileControllerRoute(profileController controllers.ProfileController) ProfileRouteController {
	return ProfileRouteController{profileController}
}

func (r *ProfileRouteController) ProfileRoute(rg *gin.RouterGroup, service services.ProfileServices) {
	router := rg.Group("/profile")

	router.GET("/me", r.profileController.FindProfile)
	router.POST("/create", r.profileController.CreateProfile)
	router.PATCH("/update", r.profileController.UpdateProfile)
	router.DELETE("/delete", r.profileController.DeleteProfile)
}
