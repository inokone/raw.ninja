package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/inokone/photostorage/web/controller"
)

func Init(v1 *gin.RouterGroup) {
	v1.GET("healthcheck", controller.Healthcheck)

	g := v1.Group("/auth")
	{
		g.POST("/login", controller.Login)
		g.GET("/logout", controller.Login)
		g.POST("/register", controller.Register)
		g.POST("/reset", controller.Reset)

	}

	g = v1.Group("/photos")
	{
		g.POST("/", controller.Upload)
		g.GET("/", controller.List)
		g.GET("/:id", controller.Get)
		g.GET("/:id/download", controller.Download)
	}
}
