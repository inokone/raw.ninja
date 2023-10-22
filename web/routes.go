package web

import (
	"github.com/gin-gonic/gin"
	"github.com/inokone/photostorage/auth"
	"github.com/inokone/photostorage/common"
	"github.com/inokone/photostorage/photo"
	"github.com/inokone/photostorage/search"
	"github.com/inokone/photostorage/statistics"
)

// Init is a function to initialize handler mapping for URLs
func Init(v1 *gin.RouterGroup, photos photo.Storer, users auth.Storer, conf common.AppConfig) {
	p := photo.NewController(photos)
	m := auth.NewJWTHandler(users, conf.Auth)
	a := auth.NewController(users, m)
	s := search.NewController(photos)
	st := statistics.NewController(photos)

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

	g = v1.Group("/statistics", m.Validate)
	{
		g.GET("/user", st.UserStatistics)
	}
}
