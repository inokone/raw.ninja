package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Health struct {
	status  string
	version string
}

// @BasePath /api/v1

// Healthcheck godoc
// @Summary Health check endpoint of the Photostorage app
// @Schemes
// @Description Returns the status and version of the application
// @Accept json
// @Produce json
// @Success 200 {Health}
// @Router /healthscheck [get]
func Healthcheck(g *gin.Context) {
	g.JSON(http.StatusOK, Health{
		status:  "ok",
		version: "0.1",
	})
}
