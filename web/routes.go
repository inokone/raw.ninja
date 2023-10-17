package web

import (
	"github.com/gin-gonic/gin"
	"github.com/inokone/photostorage/auth"
	"github.com/inokone/photostorage/common"
	"github.com/inokone/photostorage/image"
	"github.com/inokone/photostorage/photo"
	"github.com/inokone/photostorage/search"
	"gorm.io/gorm"
)

func Init(v1 *gin.RouterGroup, db *gorm.DB, ir image.Repository, conf common.AppConfig) {
	p := photo.NewController(db, ir)
	m := auth.NewJWTHandler(db, conf.Auth)
	a := auth.NewController(db, &conf.Auth)
	s := search.NewController(db, ir)

	v1.GET("healthcheck", common.Healthcheck)

	g := v1.Group("/auth")
	{
		g.POST("/login", a.Login)
		g.GET("/logout", a.Logout)
		g.POST("/signup", a.Signup)
		g.POST("/reset", a.Reset)
		g.GET("/profile", m.Validate, a.Profile)
	}

	g = v1.Group("/photos", m.Validate)
	{
		g.POST("/", p.Upload)
		g.GET("/", p.List)
		g.GET("/:id", p.Get)
		g.PUT("/:id", p.Update)
		g.DELETE("/:id", p.Delete)
		g.GET("/:id/download", p.Download)
		g.GET("/:id/thumbnail", p.Thumbnail)
	}

	g = v1.Group("/search", m.Validate)
	{
		g.GET("", s.Search)
		g.GET("/favorites", s.Favorites)
	}
}
