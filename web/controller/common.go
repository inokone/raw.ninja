package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Health struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

type StatusMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// @BasePath /api/v1

// Healthcheck godoc
// @Summary Health check endpoint of the Photostorage app
// @Schemes
// @Description Returns the status and version of the application
// @Accept json
// @Produce json
// @Success 200 {object} controller.Health
// @Router /healthcheck [get]
func Healthcheck(g *gin.Context) {
	g.JSON(http.StatusOK, Health{
		Status:  "Ok",
		Version: "0.1",
	})
}
