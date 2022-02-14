package api

import (
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var BlockDownload chan bool
var Version string

type Payload struct {
	Action string `json:"action" binding:"required,oneof=download read"`
}

func Router(logging bool) (router *gin.Engine) {
	gin.SetMode(gin.ReleaseMode)

	router = gin.New()
	router.Use(gin.Recovery())

	if logging {
		router.Use(LoggingMiddleware())
	}

	router.HandleMethodNotAllowed = true
	v1 := router.Group("/v1")

	v1.GET("/version", VersionCheck)
	v1.POST("/manage_file", ManageFile)

	BlockDownload = make(chan bool, 1)
	return
}

// Heavily inspired by https://github.com/toorop/gin-logrus
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		stop := time.Since(start)
		latency := int(math.Ceil(float64(stop.Nanoseconds()) / 1000000.0))

		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		length := c.Writer.Size()
		client := c.ClientIP()
		referer := c.Request.Referer()
		userAgent := c.Request.UserAgent()

		if len(c.Errors) > 0 {
			log.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		} else {
			msg := fmt.Sprintf(`%d %s %s %d %s "%s" "%s" %dms`, status, method, path, length, client, referer, userAgent, latency)
			log.Info(msg)
		}
	}
}

func VersionCheck(c *gin.Context) {
	c.String(http.StatusOK, Version)
}

func ManageFile(c *gin.Context) {
	var payload Payload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failure", "err": err.Error()})
		return
	}
	switch payload.Action {
	case "download":
		if assetExists() {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failure", "err": "file already downloaded"})
			return
		}
		select {
		case BlockDownload <- true:
			err := downloadFile()
			<-BlockDownload
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"status": "failure", "err": "error downloading file"})
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
