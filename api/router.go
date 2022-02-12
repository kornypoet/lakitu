package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var BlockDownload chan bool

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

	BlockDownload = make(chan bool, 1)
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
		select {
		case BlockDownload <- true:
			err := downloadAction()
			<-BlockDownload
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"status": "failure", "err": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"status": "success", "action": "download"})

		default:
			c.JSON(http.StatusTooManyRequests, gin.H{"status": "failure", "err": "file download in progress"})
		}
	case "read":
		if assetExists() {
			c.File(fmt.Sprintf("%s/%s", AssetDir, SampleFile))
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failure", "err": "file must be downloaded first"})
		}
	}
}
