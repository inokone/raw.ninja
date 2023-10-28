package web

import (
	"github.com/gin-gonic/gin"
	"github.com/inokone/photostorage/auth"
	"github.com/inokone/photostorage/auth/role"
	"github.com/inokone/photostorage/auth/user"
	"github.com/inokone/photostorage/common"
	"github.com/inokone/photostorage/mail"
	"github.com/inokone/photostorage/photo"
	"github.com/inokone/photostorage/search"
	"github.com/inokone/photostorage/stats"
)

// Init is a function to initialize handler mapping for URLs
func Init(v1 *gin.RouterGroup, photos photo.Storer, users user.Storer, roles role.Storer, auths auth.Storer, conf common.AppConfig) {
	var (
		mailer = mail.NewService(conf.Mail)
		p      = photo.NewController(photos, conf.Store)
		m      = auth.NewJWTHandler(users, conf.Auth)
		a      = auth.NewController(users, auths, m, mailer, conf.Auth)
		s      = search.NewController(photos)
		st     = stats.NewController(photos, users, conf.Store)
		u      = user.NewController(users)
		r      = role.NewController(roles)
	)

	v1.GET("healthcheck", common.Healthcheck)

	g := v1.Group("/auth")
	{
		g.POST("/login", a.Login)
		g.GET("/logout", a.Logout)
		g.POST("/signup", a.Signup)
		g.GET("/confirm", a.Confirm)
		g.POST("/resend", a.Resend)
		g.POST("/reset", u.Reset)
		g.GET("/profile", m.Validate, u.Profile)
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

	g = v1.Group("/users", m.ValidateAdmin)
	{
		g.GET("/", u.List)
		g.PATCH("/:id", u.Patch)
	}

	g = v1.Group("/roles", m.ValidateAdmin)
	{
		g.GET("/", r.List)
		g.PATCH("/:id", r.Patch)
	}

	g = v1.Group("/statistics", m.Validate)
	{
		g.GET("/user", st.UserStats)
		g.GET("/app", m.ValidateAdmin, st.AppStats)
	}
}
