package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/inokone/photostorage/auth"
	"github.com/inokone/photostorage/common"
	"github.com/inokone/photostorage/photo"
)

func Init(v1 *gin.RouterGroup) {
	v1.GET("healthcheck", common.Healthcheck)

	g := v1.Group("/auth")
	{
		g.POST("/login", auth.Login)
		g.GET("/logout", auth.Logout)
		g.POST("/register", auth.Register)
		g.POST("/reset", auth.Reset)

	}

	g = v1.Group("/photos")
	{
		g.POST("/", photo.Upload)
		g.GET("/", photo.List)
		g.GET("/:id", photo.Get)
		g.GET("/:id/download", photo.Download)
	}
}
