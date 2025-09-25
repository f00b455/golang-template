package main

import (
	"log"

	_ "github.com/f00b455/golang-template/docs" // Import generated docs
	"github.com/f00b455/golang-template/internal/config"
	"github.com/f00b455/golang-template/internal/handlers"
	"github.com/f00b455/golang-template/internal/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Golang Template API
// @version         1.0
// @description     API for Golang template project
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:3002
// @BasePath  /api

func main() {
	cfg := config.Load()

	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())

	// API routes
	api := router.Group("/api")
	{
		// Greet endpoints
		greetHandler := handlers.NewGreetHandler()
		api.GET("/greet", greetHandler.Greet)

		// RSS endpoints
		rssHandler := handlers.NewRSSHandler()
		api.GET("/rss/spiegel/latest", rssHandler.GetLatest)
		api.GET("/rss/spiegel/top5", rssHandler.GetTop5)
	}

	// Static files for terminal frontend
	router.Static("/static", "./static")
	router.StaticFile("/", "./static/terminal.html")
	router.StaticFile("/terminal", "./static/terminal.html")

	// Swagger documentation
	router.GET("/documentation/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Printf("Server starting on port %s", cfg.Port)
	log.Printf("Terminal frontend available at http://localhost:%s/", cfg.Port)
	log.Printf("Swagger documentation available at http://localhost:%s/documentation/index.html", cfg.Port)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
