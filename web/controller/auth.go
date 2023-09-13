package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/v1/auth

// Register godoc
// @Summary User registration endpoint
// @Schemes
// @Description Registers the user
// @Accept json
// @Produce json
// @Success 200 {Health}
// @Router /register [post]
func Register(g *gin.Context) {
	g.JSON(http.StatusOK, Health{
		status:  "ok",
		version: "0.1",
	})
}

// Login godoc
// @Summary User login endpoint
// @Schemes
// @Description Logs in the user, sets the necessary cookies
// @Accept json
// @Produce json
// @Success 200 {Health}
// @Router /login [post]
func Login(g *gin.Context) {
	g.JSON(http.StatusOK, Health{
		status:  "ok",
		version: "0.1",
	})
}

// Reset godoc
// @Summary Reset password endpoint
// @Schemes
// @Description Returns the status and version of the application
// @Accept json
// @Produce json
// @Success 200 {Health} Helloworld
// @Router /rest [post]
func Reset(g *gin.Context) {
	g.JSON(http.StatusOK, Health{
		status:  "ok",
		version: "0.1",
	})
}

// HealthcheckExample godoc
// @Summary Health check endpoint
// @Schemes
// @Description Returns the status and version of the application
// @Accept json
// @Produce json
// @Success 200 {Health} Helloworld
// @Router /healthscheck [get]
func Logout(g *gin.Context) {
	g.JSON(http.StatusOK, Health{
		status:  "ok",
		version: "0.1",
	})
}
