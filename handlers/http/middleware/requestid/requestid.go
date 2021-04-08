package requestid

import (
	"github.com/gin-gonic/gin"
	guuid "github.com/google/uuid"
)

// RequestID returns request id with uuid
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := guuid.New() // uuid formant
		c.Set("request_id", id.String())
		c.Next()
	}
}
