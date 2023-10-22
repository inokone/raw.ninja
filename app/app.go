package app

import (
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/inokone/photostorage/auth"
	"github.com/inokone/photostorage/common"
	"github.com/inokone/photostorage/docs"
	"github.com/inokone/photostorage/photo"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	web "github.com/inokone/photostorage/web"

	"github.com/rs/zerolog/log"
)

var (
	config *common.AppConfig
	photos photo.Storer
	users  auth.Storer
)

func init() {
	conf, err := common.LoadConfig()
	if err != nil {
		log.Error().Err(err).Msg("Failed to load application configuration.")
		os.Exit(1)
	}
	config = conf
	initLog()
	log.Info().Msg("Photostorage app starting up...")
}

// App executes the PhotoStore web application.
func App(port int) {
	var err error
	if err = initDb(config.Database); err != nil {
		log.Error().AnErr("DatabaseError", err).Msg("Failed to set up connection to database. Application spinning down.")
		os.Exit(1)
	}
	if err = initStore(config.Store); err != nil {
		log.Error().AnErr("ConfigError", err).Msg("Failed to set up image store. Application spinning down.")
		os.Exit(1)
	}

	r := gin.New()

	// Setup middleware
	r.Use(gin.Recovery())
	cc := cors.DefaultConfig()
	cc.AllowOrigins = []string{"http://127.0.0.1:3000", "http://localhost:3000"}
	cc.AllowHeaders = []string{"Authorization", "Origin", "Content-Length", "Content-Type"}
	cc.AllowCredentials = true
	r.Use(cors.New(cc))

	r.Use(web.LoggingMiddleware)
	r.MaxMultipartMemory = 8 << 20

	// Set up Swagger
	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Set up routes
	v1 := r.Group("/api/v1")

	web.Init(v1, photos, users, *config)

	r.Run(fmt.Sprintf(":%d", port))
}
