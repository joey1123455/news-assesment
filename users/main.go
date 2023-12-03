package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joey1123455/news-aggregator-service/users/config"
	"github.com/joey1123455/news-aggregator-service/users/controllers"
	docs "github.com/joey1123455/news-aggregator-service/users/docs"
	"github.com/joey1123455/news-aggregator-service/users/middleware"
	"github.com/joey1123455/news-aggregator-service/users/routes"
	"github.com/joey1123455/news-aggregator-service/users/services"
	"github.com/joey1123455/news-aggregator-service/users/utils"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	server      *gin.Engine
	ctx         context.Context
	mongoclient *mongo.Client
	redisclient *redis.Client

	userService         services.UserService
	UserController      controllers.UserController
	UserRouteController routes.UserRouteController

	authCollection      *mongo.Collection
	authService         services.AuthService
	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouteController
)

//	@title			News Aggregator user service
//	@version		1.0
//	@description	This is a user management service for the news aggregator system
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Joseph Folayan
//	@contact.email	folayanjoey@gmail.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apiKey JWT
// @in header, cookie
// @name Authorization, access_token

// @host localhost:8000
// @BasePath /api

func main() {
	config, err := config.LoadConfig(".")

	if err != nil {
		log.Fatal("Could not load config", err)
	}

	defer mongoclient.Disconnect(ctx)

	value, err := redisclient.Get(ctx, "test").Result()

	if err == redis.Nil {
		fmt.Println("key: test does not exist")
	} else if err != nil {
		utils.LogErrorToFile("open request log file", err.Error())
		panic(err)
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{config.Origin}
	corsConfig.AllowCredentials = true

	reqLogFile, err := middleware.NewFileLogger("./logs/request-logs.txt")
	if err != nil {
		utils.LogErrorToFile("open request log file", err.Error())
	}
	resLogFile, err := middleware.NewFileLogger("./logs/response-logs.txt")
	if err != nil {
		utils.LogErrorToFile("open response log file", err.Error())
	}

	server.Use(cors.New(corsConfig))
	server.Use(middleware.RequestLogger(reqLogFile))
	server.Use(middleware.ResponseLogger(resLogFile))

	docs.SwaggerInfo.BasePath = "/api"
	router := server.Group("/api")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": value})
	})

	AuthRouteController.AuthRoute(router, userService)
	UserRouteController.UserRoute(router, userService)
	// add swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	log.Fatal(server.Run(":" + config.Port))
}

func init() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	ctx = context.TODO()

	// Connect to MongoDB
	mongoconn := options.Client().ApplyURI(config.DBUri)
	mongoclient, err := mongo.Connect(ctx, mongoconn)

	if err != nil {
		utils.LogErrorToFile("connect to mongo db", err.Error())
		panic(err)
	}

	if err := mongoclient.Ping(ctx, readpref.Primary()); err != nil {
		utils.LogErrorToFile("pinging mongo db client", err.Error())
		panic(err)
	}

	fmt.Println("MongoDB successfully connected...")

	// Connect to Redis
	redisclient = redis.NewClient(&redis.Options{
		Addr: config.RedisUri,
	})

	if _, err := redisclient.Ping(ctx).Result(); err != nil {
		utils.LogErrorToFile("pinging redis client", err.Error())
		panic(err)
	}

	err = redisclient.Set(ctx, "test", "Welcome to Golang with Redis and MongoDB", 0).Err()
	if err != nil {
		utils.LogErrorToFile("setting redis db", err.Error())
		panic(err)
	}

	fmt.Println("Redis client connected successfully...")

	// Collections
	authCollection = mongoclient.Database("golang_mongodb").Collection("users")
	userService = services.NewUserServiceImpl(authCollection, ctx)
	authService = services.NewAuthService(authCollection, ctx)
	AuthController = controllers.NewAuthController(authService, userService, ctx, authCollection)
	AuthRouteController = routes.NewAuthRouteController(AuthController)

	UserController = controllers.NewUserController(userService)
	UserRouteController = routes.NewRouteUserController(UserController)

	server = gin.Default()
}
