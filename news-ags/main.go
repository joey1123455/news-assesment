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
	"github.com/joey1123455/news-aggregator-service/news-ags/config"
	"github.com/joey1123455/news-aggregator-service/news-ags/controllers"
	"github.com/joey1123455/news-aggregator-service/news-ags/routes"
	"github.com/joey1123455/news-aggregator-service/news-ags/services"
	"github.com/swaggo/swag/example/basic/docs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	server            *gin.Engine
	ctx               context.Context
	mongoclient       *mongo.Client
	redisclient       *redis.Client
	articleCollection *mongo.Collection

	scraperService services.ScrapeArticleService
	saverService   services.ArticleSaverService

	scraperController controllers.ArticleScrapperController
	saverController   controllers.ArticleSaverController

	scraperRoutesController routes.ScrapeRouteController
	saverRouteController    routes.SaveRouteController
)

//	@title			News Aggregator user service
//	@version		1.0
//	@description	This is a user management service for the news aggregator system
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Joseph Folayan
//	@contact.email	folayanjoey@gmail.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8001
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

	server.Use(cors.New(corsConfig))

	docs.SwaggerInfo.BasePath = "/api"
	router := server.Group("/api")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": value})
	})
	scraperRoutesController.ScrapeRoute(router, scraperService)
	saverRouteController.SaveRoute(router, saverService)

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

	// Collections
	articleCollection = mongoclient.Database("golang_mongodb").Collection("articles")

	// Services
	scraperService = services.NewScrapper(ctx, redisclient, Config.ApiKey)
	saverService = services.NewArticleSaver(ctx, redisclient, articleCollection)

	// Controllers
	scraperController = controllers.NewArticleScrapperController(scraperService)
	saverController = controllers.NewArticleSaverController(saverService)

	// Routes
	scraperRoutesController = routes.NewScrapeRouteController(scraperController)
	saverRouteController = routes.NewSaverRouteController(saverController)

	server = gin.Default()
}
