package main

import (
	controller "github.com/inokone/photostorage/src/web"

	docs "github.com/inokone/photostorage/src/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	{
		eg := v1.Group("/example")
		{
			eg.GET("/helloworld", controller.Helloworld)
			eg.GET("/ping", controller.Ping)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(":8080")
}
