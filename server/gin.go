// Package server handles the Gin server webhook endpoint and CRUD requests
package server

import (
	"log"

	"github.com/gin-gonic/gin"
)

func ServeGin(inspectHeaders chan<- map[string][]string, inspectBody chan<- []byte) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(gin.Recovery())

	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
