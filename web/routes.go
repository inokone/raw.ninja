package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/inokone/photostorage/auth"
	"github.com/inokone/photostorage/common"
	"github.com/inokone/photostorage/image"
	"github.com/inokone/photostorage/photo"
	"gorm.io/gorm"
)

func Init(v1 *gin.RouterGroup, db *gorm.DB, is image.Store, conf common.AppConfig) {
	p := photo.NewController(db, is)
	a := auth.NewController(db, &conf.Auth)
	m := auth.NewJWTHandler(db, conf.Auth.JWTSecret)

	v1.GET("healthcheck", common.Healthcheck)

	g := v1.Group("/auth")
	{
		g.POST("/login", a.Login)
		g.GET("/logout", a.Logout)
		g.POST("/signup", a.Signup)
		g.POST("/reset", a.Reset)
		g.GET("/validate", m.ValidateJWTToken, a.Validate)
	}

	g = v1.Group("/photos", m.ValidateJWTToken)
	{
		g.POST("/", p.Upload)
		g.GET("/", p.List)
		g.GET("/:id", p.Get)
		g.GET("/:id/download", p.Download)
	}
}
