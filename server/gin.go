// Package server handles the Gin server webhook endpoint and CRUD requests
package server

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/peterjohnbishop/friendly-octo-enigma/server/processors"
)

func ServeGin() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(gin.Recovery())

	inspectHeaders := make(chan map[string][]string, 100)
	inspectBody := make(chan []byte, 100)

	AddWebhooks(r, inspectHeaders, inspectBody)

	go processors.MapAndMergeHeaders(inspectHeaders)
	go processors.MapBody(inspectBody)

	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
