package web

import (
	"github.com/gin-gonic/gin"
	"github.com/inokone/photostorage/album"
	"github.com/inokone/photostorage/auth"
	"github.com/inokone/photostorage/auth/account"
	"github.com/inokone/photostorage/auth/role"
	"github.com/inokone/photostorage/auth/user"
	"github.com/inokone/photostorage/collection"
	"github.com/inokone/photostorage/common"
	"github.com/inokone/photostorage/image"
	"github.com/inokone/photostorage/mail"
	"github.com/inokone/photostorage/photo"
	"github.com/inokone/photostorage/search"
	"github.com/inokone/photostorage/stats"
	"github.com/inokone/photostorage/upload"
)

// Storers is a struct to collect all `Storer` entities used by the application
type Storers struct {
	Photos      photo.Storer
	Users       user.Storer
	Roles       role.Storer
	Accounts    account.Storer
	Images      image.Storer
	Collections collection.Storer
}

// Services is a struct to collect all `Service` entities used by the application
type Services struct {
	Load photo.LoadService
}

// Init is a function to initialize handler mapping for URLs
func Init(v1 *gin.RouterGroup, st Storers, se Services, c common.AppConfig) {
	var (
		mailer   = mail.NewService(c.Mail)
		uploader = photo.NewUploadService(st.Photos, st.Images, c.Store)
		loader   = photo.NewLoadService(st.Photos, st.Images, c.Store)
		p        = photo.NewController(st.Photos, st.Images, c.Store)
		m        = auth.NewJWTHandler(st.Users, c.Auth)
		a        = auth.NewController(st.Users, st.Accounts, m, c.Auth)
		ac       = account.NewController(st.Users, st.Accounts, mailer, c.Auth)
		sea      = search.NewController(st.Photos, se.Load)
		sts      = stats.NewController(st.Photos, st.Users, st.Collections, c.Store)
		u        = user.NewController(st.Users)
		r        = role.NewController(st.Roles)
		al       = album.NewController(st.Collections, loader)
		up       = upload.NewController(st.Collections, uploader, loader)
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
		g.GET("/", p.List)
		g.GET("/:id", p.Get)
		g.PUT("/:id", p.Update)
		g.DELETE("/:id", p.Delete)
		g.GET("/:id/raw", p.Raw)
		g.GET("/:id/thumbnail", p.Thumbnail)
	}

	g = v1.Group("/uploads", m.Validate)
	{
		g.POST("/", up.Upload)
		g.GET("/", up.List)
		g.GET("/:id", up.Get)
	}

	g = v1.Group("/albums", m.Validate)
	{
		g.POST("/", al.Create)
		g.PATCH("/:id", al.Patch)
		g.GET("/", al.List)
		g.GET("/:id", al.Get)
		g.DELETE("/:id", al.Delete)
	}

	g = v1.Group("/search", m.Validate)
	{
		g.GET("", sea.Search)
		g.GET("/favorites", sea.Favorites)
	}

	g = v1.Group("/users")
	{
		g.GET("/", m.ValidateAdmin, u.List)
		g.PUT("/:id", m.Validate, u.Update)
		g.PATCH("/:id", m.ValidateAdmin, u.Patch)
		g.PUT("/:id/enabled", m.ValidateAdmin, u.SetEnabled)
	}

	g = v1.Group("/roles", m.ValidateAdmin)
	{
		g.GET("/", r.List)
		g.PUT("/:id", r.Update)
	}

	g = v1.Group("/statistics", m.Validate)
	{
		g.GET("/user", sts.UserStats)
		g.GET("/app", m.ValidateAdmin, sts.AppStats)
	}
}
