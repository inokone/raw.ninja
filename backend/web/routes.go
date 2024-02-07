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
	"github.com/inokone/photostorage/onetime"
	"github.com/inokone/photostorage/photo"
	"github.com/inokone/photostorage/ruleset"
	"github.com/inokone/photostorage/ruleset/rule"
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
	Rules       rule.Storer
	RuleSets    ruleset.Storer
	OneTime     onetime.Storer
}

// Services is a struct to collect all `Service` entities used by the application
type Services struct {
	Load photo.LoadService
}

// InitPrivate is a function to initialize handler mapping for URLs protected with CORS
func InitPrivate(private *gin.RouterGroup, st Storers, se Services, c *common.AppConfig) {
	var (
		mailer   = mail.NewService(c.Mail)
		colls    = collection.NewService(st.Collections)
		uploader = photo.NewUploadService(st.Photos, st.Images, c.Store)
		loader   = photo.NewLoadService(st.Photos, st.Images, c.Store)
		p        = photo.NewController(st.Photos, st.Images, c.Store)
		m        = auth.NewJWTHandler(st.Users, c.Auth)
		a        = auth.NewController(st.Users, st.Accounts, m, c.Auth)
		ac       = account.NewController(st.Users, st.Accounts, mailer, c.Auth)
		sea      = search.NewController(st.Photos, se.Load, colls)
		sts      = stats.NewController(st.Photos, st.Users, st.Collections, c.Store)
		u        = user.NewController(st.Users)
		r        = role.NewController(st.Roles)
		al       = album.NewController(st.Collections, loader, colls)
		up       = upload.NewController(st.Collections, uploader, loader, colls)
		rs       = ruleset.NewController(st.RuleSets, st.Rules)
		ru       = rule.NewController(st.Rules)
		ot       = onetime.NewController(st.OneTime, st.Images)
	)

	private.GET("healthcheck", common.Healthcheck)

	g := private.Group("/auth")
	{
		g.POST("/login", a.Login)
		g.GET("/logout", a.Logout)
	}

	g = private.Group("/account")
	{
		g.POST("/signup", ac.Signup)
		g.GET("/confirm", ac.Confirm)
		g.PUT("/resend", ac.ResendConfirmation)
		g.PUT("/recover", ac.Recover)
		g.PUT("/password/reset", ac.ResetPassword)
		g.PUT("/password/change", m.Validate, ac.ChangePassword)
		g.GET("/profile", m.Validate, u.Profile)
	}

	g = private.Group("/photos", m.Validate)
	{
		g.GET("/", p.List)
		g.GET("/:id", p.Get)
		g.PUT("/:id", p.Update)
		g.DELETE("/:id", p.Delete)
		g.GET("/:id/raw", p.Raw)
		g.GET("/:id/thumbnail", p.Thumbnail)
	}

	g = private.Group("/onetime", m.Validate)
	{
		g.POST("/", ot.Create)
	}

	g = private.Group("/uploads", m.Validate)
	{
		g.POST("/", up.Upload)
		g.GET("/", up.List)
		g.GET("/:id", up.Get)
	}

	g = private.Group("/albums", m.Validate)
	{
		g.POST("/", al.Create)
		g.PATCH("/:id", al.Patch)
		g.GET("/", al.List)
		g.GET("/:id", al.Get)
		g.DELETE("/:id", al.Delete)
	}

	g = private.Group("/search", m.Validate)
	{
		g.GET("", sea.Search)
		g.GET("/favorites", sea.Favorites)
	}

	g = private.Group("/users")
	{
		g.GET("/", m.ValidateAdmin, u.List)
		g.PUT("/:id", m.Validate, u.Update)
		g.PATCH("/:id", m.ValidateAdmin, u.Patch)
		g.PUT("/:id/enabled", m.ValidateAdmin, u.SetEnabled)
	}

	g = private.Group("/roles", m.ValidateAdmin)
	{
		g.GET("/", r.List)
		g.PUT("/:id", r.Update)
	}

	g = private.Group("/rules", m.Validate)
	{
		g.POST("/", ru.Create)
		g.GET("/", ru.List)
		g.GET("/constants", ru.Constants)
		g.GET("/:id", ru.Get)
	}

	g = private.Group("/rulesets", m.Validate)
	{
		g.POST("/", rs.Create)
		g.GET("/", rs.List)
		g.GET("/:id", rs.Get)
		g.PUT("/:id", rs.Update)
		g.DELETE("/:id", rs.Delete)
	}

	g = private.Group("/statistics", m.Validate)
	{
		g.GET("/user", sts.UserStats)
		g.GET("/users", sts.Users)
		g.GET("/app", m.ValidateAdmin, sts.AppStats)
	}
}

// InitPublic is a function to initialize handler mapping for URLs not protected with CORS
func InitPublic(public *gin.RouterGroup, st Storers, c *common.AppConfig) {
	ot := onetime.NewController(st.OneTime, st.Images)
	m := auth.NewJWTHandler(st.Users, c.Auth)
	gt := auth.NewGoogleController(*c.Auth, st.Users, m)
	ft := auth.NewFacebookController(*c.Auth, st.Users, m)

	g := public.Group("/onetime")
	{
		g.GET("/raw/:id", ot.Raw)
	}

	g = public.Group("/auth")
	{
		g.GET("/google/redirect", gt.Redirect)
		g.GET("/google", gt.Login)
		g.GET("/facebook/redirect", ft.Redirect)
		g.GET("/facebook", ft.Login)
	}
}
