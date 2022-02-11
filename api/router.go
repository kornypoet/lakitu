package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Payload struct {
	Action string `json:"action" binding:"required,oneof=download read"`
}

func Router(debug bool) *gin.Engine {
	gin.DisableConsoleColor()

	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.HandleMethodNotAllowed = true
	v1 := router.Group("/v1")

	v1.GET("/ping", Ping)
	v1.POST("/manage_file", ManageFile)

	return router
}

func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func ManageFile(c *gin.Context) {
	var payload Payload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failure", "err": err.Error()})
		return
	}
	switch payload.Action {
	case "download":
		// stub
		c.JSON(http.StatusOK, gin.H{"status": "success", "action": "download"})
	case "read":
		// stub
		c.JSON(http.StatusOK, gin.H{"status": "success", "action": "read"})
	}
}
