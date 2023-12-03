package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joey1123455/news-aggregator-service/users/config"
	"github.com/joey1123455/news-aggregator-service/users/models"
	"github.com/joey1123455/news-aggregator-service/users/services"
	"github.com/joey1123455/news-aggregator-service/users/utils"
	"github.com/thanhpk/randstr"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthController struct {
	authService services.AuthService
	userService services.UserService
	ctx         context.Context
	collection  *mongo.Collection
}

func NewAuthController(authService services.AuthService, userService services.UserService, ctx context.Context, collection *mongo.Collection) AuthController {
	return AuthController{authService, userService, ctx, collection}
}

// @Summary Sign up user
// @Description Sign up a new user with the provided information.
// @Produce json
// @Accept json
// @Param user body models.SignUpInput true "User information for sign up"
// @Success 201 {object} gin.H{"status": "success", "message": "We sent an email with a verification code to user@example.com"} "User signed up successfully"
// @Failure 400 {object} gin.H{"status": "fail", "message": "Invalid request payload or parameters"} "Invalid request payload or parameters"
// @Failure 400 {object} gin.H{"status": "fail", "message": "Passwords do not match"} "Passwords do not match"
// @Failure 409 {object} gin.H{"status": "error", "message": "Email already exists"} "Email already exists"
// @Failure 502 {object} gin.H{"status": "error", "message": "Error while signing up new user"} "Error while signing up new user"
// @Router /sign-up [post]
func (ac *AuthController) SignUpUser(ctx *gin.Context) {
	var user *models.SignUpInput

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if user.Password != user.PasswordConfirm {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Passwords do not match"})
		return
	}

	newUser, err := ac.authService.SignUpUser(user)

	if err != nil {
		utils.LogErrorToFile("while signing up new user", err.Error())
		if strings.Contains(err.Error(), "email already exist") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "error", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		return
	}

	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load config", err)
	}

	// Generate Verification Code
	code := randstr.String(20)

	verificationCode := utils.Encode(code)

	// Update User in Database
	ac.userService.UpdateUserById(newUser.ID.Hex(), "verificationCode", verificationCode)

	var firstName = newUser.Name

	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	// 👇 Send Email
	emailData := utils.EmailData{
		URL:       config.Origin + "/verifyemail/" + code,
		FirstName: firstName,
		Subject:   "Your account verification code",
	}

	utils.SendEmail(newUser, &emailData, "verificationCode.html")

	message := "We sent an email with a verification code to " + user.Email
	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": message})
}

// @Summary Sign in user
// @Description Sign in a user using their email and password.
// @Produce json
// @Accept json
// @Param credentials body models.SignInInput true "User credentials for sign in"
// @Success 200 {object} gin.H{"status": "success", "access_token": "access_token"} "User signed in successfully"
// @Failure 400 {object} gin.H{"status": "fail", "message": "Invalid email or password"} "Invalid email or password"
// @Failure 401 {object} gin.H{"status": "fail", "message": "You are not verified, please verify your email to login"} "User not verified"
// @Router /auth/sign-in [post]
func (ac *AuthController) SignInUser(ctx *gin.Context) {
	var credentials *models.SignInInput

	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	user, err := ac.userService.FindUserByEmail(credentials.Email)
	if err != nil {
		utils.LogErrorToFile("while signing in a user", err.Error())
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or password"})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if !user.Verified {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are not verified, please verify your email to login"})
		return
	}

	if err := utils.VerifyPassword(user.Password, credentials.Password); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}

	config, _ := config.LoadConfig(".")

	// Generate Tokens
	access_token, err := utils.CreateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenPrivateKey)
	if err != nil {
		utils.LogErrorToFile("while generating access token for user sign in", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	refresh_token, err := utils.CreateToken(config.RefreshTokenExpiresIn, user.ID, config.RefreshTokenPrivateKey)
	if err != nil {
		utils.LogErrorToFile("while generating refresh token for user sing in", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("access_token", access_token, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refresh_token, config.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": access_token})
}

// @Summary Refresh access token
// @Description Refresh the user's access token using the provided refresh token.
// @Produce json
// @Success 200 {object} gin.H{"status": "success", "access_token": "new_access_token"} "Access token refreshed successfully"
// @Failure 403 {object} gin.H{"status": "fail", "message": "Could not refresh access token"} "Could not refresh access token"
// @Router /auth/refresh-access-token [get]
func (ac *AuthController) RefreshAccessToken(ctx *gin.Context) {
	message := "could not refresh access token"

	cookie, err := ctx.Cookie("refresh_token")

	if err != nil {
		utils.LogErrorToFile("while reading refresh token for user sign in", err.Error())
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": message})
		return
	}

	config, _ := config.LoadConfig(".")

	sub, err := utils.ValidateToken(cookie, config.RefreshTokenPublicKey)
	if err != nil {
		utils.LogErrorToFile("while validating refresh token for user sign in", err.Error())
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	user, err := ac.userService.FindUserById(fmt.Sprint(sub))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token no logger exists"})
		return
	}

	access_token, err := utils.CreateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenPrivateKey)
	if err != nil {
		utils.LogErrorToFile("while creating access token for user sign in, from token refresh", err.Error())
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("access_token", access_token, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": access_token})
}

// @Summary Logout a user
// @Description Logout a user by clearing the access_token, refresh_token, and logged_in cookies.
// @Produce json
// @Success 200 {object} gin.H{"status": "success"} "Successfully logged out"
// @Router /auth/logout [post]
func (ac *AuthController) LogoutUser(ctx *gin.Context) {
	ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

// @Summary Verify email address
// @Description Verify email address using the provided verification code.
// @Produce json
// @Param verificationCode path string true "Verification code for email verification"
// @Success 200 {object} gin.H{"status": "success", "message": "Email verified successfully"} "Email verified successfully"
// @Failure 403 {object} gin.H{"status": "success", "message": "Could not verify email address"} "Could not verify email address"
// @Failure 403 {object} gin.H{"status": "success", "message": "Error while verifying email"} "Error while verifying email"
// @Router /auth/verify-email/{verificationCode} [post]
func (ac *AuthController) VerifyEmail(ctx *gin.Context) {

	code := ctx.Params.ByName("verificationCode")
	verificationCode := utils.Encode(code)

	query := bson.D{{Key: "verificationCode", Value: verificationCode}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "verified", Value: true}}}, {Key: "$unset", Value: bson.D{{Key: "verificationCode", Value: ""}}}}
	result, err := ac.collection.UpdateOne(ac.ctx, query, update)
	if err != nil {
		utils.LogErrorToFile("while verifying email", err.Error())
		ctx.JSON(http.StatusForbidden, gin.H{"status": "success", "message": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		ctx.JSON(http.StatusForbidden, gin.H{"status": "success", "message": "Could not verify email address"})
		return
	}

	fmt.Println(result)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Email verified successfully"})

}

// @Summary Reset password
// @Description Reset user password using the provided reset token and new password.
// @Produce json
// @Accept json
// @Param resetToken path string true "Reset token for password reset"
// @Param userCredential body models.ResetPasswordInput true "User credentials for password reset"
// @Success 200 {object} gin.H{"status": "success", "message": "Password data updated successfully"} "Password data updated successfully"
// @Failure 400 {object} gin.H{"status": "fail", "message": "Invalid request payload or parameters"} "Invalid request payload or parameters"
// @Failure 400 {object} gin.H{"status": "fail", "message": "Passwords do not match"} "Passwords do not match"
// @Failure 400 {object} gin.H{"status": "fail", "message": "Token is invalid or has expired"} "Token is invalid or has expired"
// @Failure 403 {object} gin.H{"status": "fail", "message": "Error while resetting password"} "Error while resetting password"
// @Router /auth/reset-password/{resetToken} [post]
func (ac *AuthController) ForgotPassword(ctx *gin.Context) {
	var userCredential *models.ForgotPasswordInput

	if err := ctx.ShouldBindJSON(&userCredential); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	message := "You will receive a reset email if user with that email exist"

	user, err := ac.userService.FindUserByEmail(userCredential.Email)
	if err != nil {
		utils.LogErrorToFile("while resetting password", err.Error())
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusOK, gin.H{"status": "fail", "message": message})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if !user.Verified {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Account not verified"})
		return
	}

	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load config", err)
	}

	// Generate Verification Code
	resetToken := randstr.String(20)

	passwordResetToken := utils.Encode(resetToken)

	// Update User in Database
	query := bson.D{{Key: "email", Value: strings.ToLower(userCredential.Email)}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "passwordResetToken", Value: passwordResetToken}, {Key: "passwordResetAt", Value: time.Now().Add(time.Minute * 15)}}}}
	result, err := ac.collection.UpdateOne(ac.ctx, query, update)

	if result.MatchedCount == 0 {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "success", "message": "There was an error sending email"})
		return
	}

	if err != nil {
		utils.LogErrorToFile("while resetting password", err.Error())
		ctx.JSON(http.StatusForbidden, gin.H{"status": "success", "message": err.Error()})
		return
	}
	var firstName = user.Name

	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	// ? Send Email
	emailData := utils.EmailData{
		URL:       config.Origin + "/resetpassword/" + resetToken,
		FirstName: firstName,
		Subject:   "Your password reset token (valid for 10min)",
	}

	err = utils.SendEmail(user, &emailData, "resetPassword.html")
	if err != nil {
		utils.LogErrorToFile("while sending email for password reset", err.Error())
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "success", "message": "There was an error sending email"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
}

// @Summary Reset password
// @Description Reset user password using the provided reset token and new password.
// @Produce json
// @Accept json
// @Param resetToken path string true "Reset token for password reset"
// @Param userCredential body models.ResetPasswordInput true "User credentials for password reset"
// @Success 200 {object} gin.H{"status": "success", "message": "Password data updated successfully"} "Password data updated successfully"
// @Failure 400 {object} gin.H{"status": "fail", "message": "Invalid request payload or parameters"} "Invalid request payload or parameters"
// @Failure 400 {object} gin.H{"status": "fail", "message": "Passwords do not match"} "Passwords do not match"
// @Failure 400 {object} gin.H{"status": "fail", "message": "Token is invalid or has expired"} "Token is invalid or has expired"
// @Failure 403 {object} gin.H{"status": "fail", "message": "Error while resetting password"} "Error while resetting password"
// @Router /auth/reset-password/{resetToken} [post]
func (ac *AuthController) ResetPassword(ctx *gin.Context) {
	resetToken := ctx.Params.ByName("resetToken")
	var userCredential *models.ResetPasswordInput

	if err := ctx.ShouldBindJSON(&userCredential); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if userCredential.Password != userCredential.PasswordConfirm {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Passwords do not match"})
		return
	}

	hashedPassword, _ := utils.HashPassword(userCredential.Password)

	passwordResetToken := utils.Encode(resetToken)

	// Update User in Database
	query := bson.D{{Key: "passwordResetToken", Value: passwordResetToken}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "password", Value: hashedPassword}}}, {Key: "$unset", Value: bson.D{{Key: "passwordResetToken", Value: ""}, {Key: "passwordResetAt", Value: ""}}}}
	result, err := ac.collection.UpdateOne(ac.ctx, query, update)

	if result.MatchedCount == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "success", "message": "Token is invalid or has expired"})
		return
	}

	if err != nil {
		utils.LogErrorToFile("while resetting password", err.Error())
		ctx.JSON(http.StatusForbidden, gin.H{"status": "success", "message": err.Error()})
		return
	}

	ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Password data updated successfully"})
}
