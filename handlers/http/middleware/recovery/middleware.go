package recovery

import (
	"bytes"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"simpleBackend/log"
)

type response map[string]interface{}

// PanicRecovery handles panic and responses message as json format
// HINT: just keeping simple for error message
func PanicRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {

		// in go-gin, the body is just picked up once
		// so, if we need to log request body, we need copy it!
		var buf bytes.Buffer
		var body []byte
		if c.Request.Body != nil {
			tee := io.TeeReader(c.Request.Body, &buf)
			body, _ = ioutil.ReadAll(tee)
			c.Request.Body = ioutil.NopCloser(&buf)
		}

		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							c.Error(err.(error))
							c.Abort()
							return
						}
					}
				}

				log.Logger.Error("panic error", zap.String("http-request-body", string(body)))

				switch v := err.(type) {
				case APIError:
					c.AbortWithStatusJSON(v.Status, v.response())
				case *APIError:
					c.AbortWithStatusJSON(v.Status, v.response())
				case string:
					apiErr := APIError{Msg: v, Code: InternalError}
					c.AbortWithStatusJSON(http.StatusInternalServerError, apiErr.response())
				default:
					c.AbortWithStatusJSON(http.StatusInternalServerError, response{
						"code":  InternalError,
						"error": v,
					})
				}

			}
		}()
		c.Next()
	}
}
