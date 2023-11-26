package app

import (
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/inokone/photostorage/common"
	"github.com/inokone/photostorage/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	web "github.com/inokone/photostorage/web"

	"github.com/rs/zerolog/log"
)

var (
	config   *common.AppConfig
	storers  web.Storers
	services web.Services
)

// App executes the RAW.Ninja web application.
func App(path string) {
	var err error
	if err = initConf(path); err != nil {
		log.Error().Err(err).Msg("Failed to load application configuration.")
		os.Exit(1)
	}
	if err = initDb(config.Database); err != nil {
		log.Error().AnErr("DatabaseError", err).Msg("Failed to set up connection to database. Application spinning down.")
		os.Exit(1)
	}
	initStorers(config.Store)
	initServices(config.Store, storers)

	r := gin.New()

	// Setup middleware
	r.Use(gin.Recovery())
	cc := cors.DefaultConfig()
	cc.AllowOrigins = []string{"https://raw.ninja", "https://rawninja.net", config.Auth.FrontendRoot}
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

	web.Init(v1, storers, services, *config)

	p := fmt.Sprintf("0.0.0.0:%d", config.Web.Port)
	if len(config.Auth.TLSCert) > 0 {
		err = r.RunTLS(p, config.Auth.TLSCert, config.Auth.TLSKey)
		if err != nil {
			log.Err(err).Msg("Failed to start the application")
		}
	} else {
		err = r.Run(p)
		if err != nil {
			log.Err(err).Msg("Failed to start the application")
		}
	}
}
