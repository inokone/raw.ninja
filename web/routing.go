package controller

import (
	"github.com/gin-gonic/gin"
)

func Init(v1 *gin.RouterGroup) {
	{
		eg := v1.Group("/example")
		{
			eg.GET("/helloworld", Helloworld)
			eg.GET("/ping", Ping)
		}
	}
}
