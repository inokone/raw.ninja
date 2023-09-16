package app

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/inokone/photostorage/common"
	"github.com/inokone/photostorage/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	routes "github.com/inokone/photostorage/web"
	"gorm.io/gorm"

	"github.com/inokone/photostorage/image"
)

var config *common.AppConfig
var DB *gorm.DB
var IS *image.Store

func init() {
	conf, err := common.LoadConfig()
	if err != nil {
		log.Fatal("Could not load application configuration", err)
	}
	config = conf
}

func App(port int) {
	err := initDb(config.Database)
	if err != nil {
		log.Fatal("Could not set up connection to database. Application spinning down.")
		os.Exit(1)
	}

	err = initStore(config.Store)
	if err != nil {
		log.Fatal("Could not set up image store. Application spinning down.")
		os.Exit(1)
	}

	r := gin.Default()

	// Set up Swagger
	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Set up routes
	v1 := r.Group("/api/v1")
	routes.Init(v1, DB, *IS, *config)

	r.Run(fmt.Sprintf(":%d", port))
}
