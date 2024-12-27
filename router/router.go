package router

import (
	"github.com/gin-gonic/gin"
)

func Run() {
	r := gin.Default()

	r.GET("/appstacks", handleAppStacksGET)

	r.POST("/appstacks", handleAppStacksPOST)

	r.POST("/webhook/event/:appname", handleEventWebhookPOST)

	r.Run(":8080")
}
