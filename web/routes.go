package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/inokone/photostorage/auth"
	"github.com/inokone/photostorage/common"
	"github.com/inokone/photostorage/image"
	"github.com/inokone/photostorage/photo"
	"gorm.io/gorm"
)

func Init(v1 *gin.RouterGroup, db *gorm.DB, is image.Store) {
	p := photo.NewController(db, is)

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
		g.POST("/", p.Upload)
		g.GET("/", p.List)
		g.GET("/:id", p.Get)
		g.GET("/:id/download", p.Download)
	}
}
