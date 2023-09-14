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
// @Success 200
// @Router /register [post]
func Register(g *gin.Context) {
	g.JSON(http.StatusNotImplemented, StatusMessage{
		Code:    501,
		Message: "Functionality has not been implemented yet!",
	})
}

// Login godoc
// @Summary User login endpoint
// @Schemes
// @Description Logs in the user, sets the necessary cookies
// @Accept json
// @Produce json
// @Success 200
// @Router /login [post]
func Login(g *gin.Context) {
	g.JSON(http.StatusNotImplemented, StatusMessage{
		Code:    501,
		Message: "Functionality has not been implemented yet!",
	})
}

// Reset godoc
// @Summary Reset password endpoint
// @Schemes
// @Description Returns the status and version of the application
// @Accept json
// @Produce json
// @Success 200
// @Router /reset [post]
func Reset(g *gin.Context) {
	g.JSON(http.StatusNotImplemented, StatusMessage{
		Code:    501,
		Message: "Functionality has not been implemented yet!",
	})
}

// Logout godoc
// @Summary Logout endpoint
// @Schemes
// @Description Logs out of the application
// @Accept json
// @Produce json
// @Success 200
// @Router /logout [get]
func Logout(g *gin.Context) {
	g.JSON(http.StatusNotImplemented, StatusMessage{
		Code:    501,
		Message: "Functionality has not been implemented yet!",
	})
}
