package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Router(debug bool) *gin.Engine {
	gin.DisableConsoleColor()

	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.HandleMethodNotAllowed = true
	v1 := router.Group("/v1")

	v1.GET("/ping", Ping)

	return router
}

func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
