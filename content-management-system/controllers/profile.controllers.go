package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joey1123455/news-aggregator-service/content-management-system/models"
	"github.com/joey1123455/news-aggregator-service/content-management-system/services"
)

type ProfileController struct {
	service services.ProfileServices
}

func NewProfileController(postService services.ProfileServices) ProfileController {
	return ProfileController{postService}
}

// @Summary Create User Profile
// @Description Create a users profile.
// @Security ApiKeyAuth
// @Produce json
// @Accept json
// @Param user body Profile true "User information for profile creation"
// @Success 201 {object} Profile
// @Failure 409 {object} string "profile already exists"
// @Failure 502 {object} string "error message"
// @Router /profile/create [post]
func (pc *ProfileController) CreateProfile(ctx *gin.Context) {
	var data *models.CreateUser
	id := ctx.MustGet("currentUserId").(string)

	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	newProfile, err := pc.service.CreateUser(id, data)

	if err != nil {
		if strings.Contains(err.Error(), "profile already exists") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": err.Error()})
			return
		}

		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	profile := models.FilteredResponse(*newProfile)

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "profile": profile})
}

// @Summary Retrieve User Profile
// @Description Retrieves a users profile.
// @Security ApiKeyAuth
// @Produce json
// @Accept json
// @Success 200 {object} Profile
// @Failure 404 {object} string "User not found"
// @Failure 502 {object} string "Error while signing up new user"
// @Router /profile/me [get]
func (pc *ProfileController) FindProfile(ctx *gin.Context) {
	id := ctx.MustGet("currentUserId").(string)

	profile, err := pc.service.FindUser(id)

	if err != nil {
		if strings.Contains(err.Error(), "Id exists") {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": models.FilteredResponse(*profile)})
}

// @Summary Update User Profile
// @Description Update a users profile.
// @Security ApiKeyAuth
// @Produce json
// @Accept json
// @Param user body models.UpdateUser true "User information for profile update"
// @Success 200 {object} Profile
// @Failure 404 {object} string "User not found"
// @Failure 502 {object} string "Error while signing up new user"
// @Router /profile/update [patch]
func (pc *ProfileController) UpdateProfile(ctx *gin.Context) {
	id := ctx.MustGet("currentUserId").(string)

	var data *models.UpdateUser
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	updatedProfile, err := pc.service.UpdateUser(id, data)
	if err != nil {
		if strings.Contains(err.Error(), "Id exists") {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	profile := models.FilteredResponse(*updatedProfile)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "profile": profile})
}

// @Summary Delete User Profile
// @Description Delete a users profile.
// @Security ApiKeyAuth
// @Produce json
// @Accept json
// @Success 204
// @Failure 404 {object} string "User not found"
// @Failure 502 {object} string "Error while signing up new user"
// @Router /profile/delete [delete]
func (pc *ProfileController) DeleteProfile(ctx *gin.Context) {
	id := ctx.MustGet("currentUserId").(string)

	err := pc.service.DeleteUser(id)

	if err != nil {
		if strings.Contains(err.Error(), "Id exists") {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
