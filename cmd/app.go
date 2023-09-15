package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/inokone/photostorage/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	routes "github.com/inokone/photostorage/web"
)

func App(port int) {
	r := gin.Default()

	// Set up Swagger
	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Set up routes
	v1 := r.Group("/api/v1")
	routes.Init(v1, DB, *IS)

	r.Run(fmt.Sprintf(":%d", port))
}
