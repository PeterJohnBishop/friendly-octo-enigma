package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AddWebhooks imports webhook endpoints into the Gin router
func AddWebhooks(r *gin.Engine, inspectHeaders chan map[string][]string, inspectBody chan []byte) {
	v1 := r.Group("/v1", nil)

	// sends request headers and body to separate channels for processing
	v1.POST("/webhook/inspect", func(c *gin.Context) {
		bodyBytes, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
			return
		}

		for key, values := range c.Request.Header {
			h := map[string][]string{
				key: values,
			}
			// send headers to the inspectHeaders channel
			inspectHeaders <- h
		}

		// send body []byte to inspectBody channel
		inspectBody <- bodyBytes

		c.JSON(http.StatusOK, gin.H{
			"message": "Payload recieved.",
		})
	})
}
