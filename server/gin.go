// Package server handles the Gin server webhook endpoint and CRUD requests
package server

import (
	"log"

	"github.com/gin-gonic/gin"
)

func ServeGin() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(gin.Recovery())

	inspectHeaders := make(chan map[string][]string)
	inspectBody := make(chan []byte)

	AddWebhooks(r, inspectHeaders, inspectBody)

	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
