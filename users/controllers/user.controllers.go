package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joey1123455/news-aggregator-service/users/models"
	"github.com/joey1123455/news-aggregator-service/users/services"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) UserController {
	return UserController{userService}
}

// @Summary Get current user details
// @Description Get details of the currently authenticated user.
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.UserResponse "Current user details retrieved successfully"
// @Failure 401 {object} string "Unauthorized access"
// @Router /user/me [get]
func (uc *UserController) GetMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(*models.DBResponse)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": models.FilteredResponse(currentUser)}})
}
