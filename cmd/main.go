package main

import (
	"log"

	"github.com/inokone/photostorage/common"
	"github.com/inokone/photostorage/docs"
	controller "github.com/inokone/photostorage/web"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var config *common.AppConfig

func init() {
	conf, err := common.LoadConfig()
	if err != nil {
		log.Fatal("Could not load application configuration", err)
	}
	config = conf
}

func main() {
	common.InitDb(config.Database)

	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	controller.Init(v1)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(":8080")
}
