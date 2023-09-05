package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/v1

// HelloWorldExample godoc
// @Summary hello world example
// @Schemes
// @Description do hello world
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /example/helloworld [get]
func Helloworld(g *gin.Context) {
	g.JSON(http.StatusOK, "helloworld")
}

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping-pong
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /example/ping [get]
func Ping(g *gin.Context) {
	g.JSON(http.StatusOK, gin.H{"message": "pong"})
}
