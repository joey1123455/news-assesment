package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/joey1123455/news-aggregator-service/content-management-system/config"
	"github.com/joey1123455/news-aggregator-service/content-management-system/controllers"
	docs "github.com/joey1123455/news-aggregator-service/content-management-system/docs"
	"github.com/joey1123455/news-aggregator-service/content-management-system/routes"
	"github.com/joey1123455/news-aggregator-service/content-management-system/services"
	"github.com/joey1123455/news-aggregator-service/content-management-system/utils"

	middleware "github.com/joey1123455/news-aggregator-service/content-management-system/middlewares"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	server            *gin.Engine
	ctx               context.Context
	mongoclient       *mongo.Client
	redisclient       *redis.Client
	profileCollection *mongo.Collection
	newsCollection    *mongo.Collection

	profileService services.ProfileServices
	newsService    services.ArticleServices

	newsController    controllers.NewsController
	profileController controllers.ProfileController

	newsRouter    routes.NewsRouteController
	profileRouter routes.ProfileRouteController
)

//	@title			News aggregator content management service
//	@version		1.0
//	@description	This is a content management service for the news aggregator system
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Joseph Folayan
//	@contact.email	folayanjoey@gmail.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apiKey JWT
// @in header, cookie
// @name Authorization, access_token

// @host localhost:8002
// @BasePath /api

func main() {
	config, err := config.LoadConfig(".")

	if err != nil {
		log.Fatal("Could not load config", err)
	}

	defer mongoclient.Disconnect(ctx)

	value, err := redisclient.Get("test").Result()

	if err == redis.Nil {
		fmt.Println("key: test does not exist")
	} else if err != nil {
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

	newsRouter.NewsRoute(router, newsService)
	profileRouter.ProfileRoute(router, profileService)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	log.Fatal(server.Run(":" + config.Port))
}

func init() {
	Config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	ctx = context.TODO()

	// Connect to MongoDB
	mongoconn := options.Client().ApplyURI(Config.DBUri)
	mongoclient, err := mongo.Connect(ctx, mongoconn)

	if err != nil {
		// utils.LogErrorToFile("connect to mongo db", err.Error())
		panic(err)
	}

	if err := mongoclient.Ping(ctx, readpref.Primary()); err != nil {
		// utils.LogErrorToFile("pinging mongo db client", err.Error())
		panic(err)
	}

	fmt.Println("MongoDB successfully connected...")

	// Connect to Redis
	redisclient = redis.NewClient(&redis.Options{
		Addr: Config.RedisUri,
	})

	if _, err := redisclient.Ping().Result(); err != nil {
		// utils.LogErrorToFile("pinging redis client", err.Error())
		panic(err)
	}

	err = redisclient.Set("test", "Connected to Redis and MongoDB", time.Hour).Err()
	if err != nil {
		// utils.LogErrorToFile("setting redis db", err.Error())
		panic(err)
	}

	fmt.Println("Redis client connected successfully...")

	newsCollection = mongoclient.Database("golang_mongodb").Collection("articles")
	profileCollection = mongoclient.Database("golang_mongodb").Collection("profiles")

	newsService = services.NewArticleService(ctx, newsCollection)
	profileService = services.NewProfileService(ctx, profileCollection)

	newsController = controllers.NewNewsController(newsService)
	profileController = controllers.NewProfileController(profileService)

	newsRouter = routes.NewNewsControllerRoute(newsController)
	profileRouter = routes.NewprofileControllerRoute(profileController)

	server = gin.Default()
}
