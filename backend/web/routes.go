package web

import (
	"github.com/gin-gonic/gin"
	"github.com/inokone/photostorage/auth"
	"github.com/inokone/photostorage/auth/account"
	"github.com/inokone/photostorage/auth/role"
	"github.com/inokone/photostorage/auth/user"
	"github.com/inokone/photostorage/common"
	"github.com/inokone/photostorage/mail"
	"github.com/inokone/photostorage/photo"
	"github.com/inokone/photostorage/search"
	"github.com/inokone/photostorage/stats"
)

// Storers is a struct to collect all `Storer` entities used by the application
type Storers struct {
	Photos   photo.Storer
	Users    user.Storer
	Roles    role.Storer
	Accounts account.Storer
}

// Init is a function to initialize handler mapping for URLs
func Init(v1 *gin.RouterGroup, s Storers, c common.AppConfig) {
	var (
		mailer = mail.NewService(c.Mail)
		p      = photo.NewController(s.Photos, c.Store)
		m      = auth.NewJWTHandler(s.Users, c.Auth)
		a      = auth.NewController(s.Users, s.Accounts, m)
		ac     = account.NewController(s.Users, s.Accounts, mailer, c.Auth)
		se     = search.NewController(s.Photos)
		st     = stats.NewController(s.Photos, s.Users, c.Store)
		u      = user.NewController(s.Users)
		r      = role.NewController(s.Roles)
	)

	v1.GET("healthcheck", common.Healthcheck)

	g := v1.Group("/auth")
	{
		g.POST("/login", a.Login)
		g.GET("/logout", a.Logout)
	}

	g = v1.Group("/account")
	{
		g.POST("/signup", ac.Signup)
		g.GET("/confirm", ac.Confirm)
		g.PUT("/resend", ac.ResendConfirmation)
		g.PUT("/recover", ac.Recover)
		g.PUT("/password/reset", ac.ResetPassword)
		g.PUT("/password/change", m.Validate, ac.ChangePassword)
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
		g.GET("", se.Search)
		g.GET("/favorites", se.Favorites)
	}

	g = v1.Group("/users")
	{
		g.GET("/", m.ValidateAdmin, u.List)
		g.PUT("/:id", m.Validate, u.Update)
		g.PATCH("/:id", m.ValidateAdmin, u.Patch)
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
