package httphandler

import (
	"net/http"
	"simpleBackend/handlers/http/middleware/cors"
	"simpleBackend/handlers/http/middleware/recovery"
	"simpleBackend/handlers/http/middleware/requestid"

	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
)

// import "github.com/pkg/errors"

// Option is option for handler
type Option func(*Handler) error

// Handler handles the http service
type Handler struct {
	Mode string
	// TODO db
	// nasn api key
}

// New returns http handler
func New(mode string, opts ...Option) (*Handler, error) {
	h := &Handler{Mode: mode}

	for _, opt := range opts {
		if err := opt(h); err != nil {
			return nil, err
		}
	}

	return h, nil
}

func newRouter(mode string) *gin.Engine {
	switch mode {
	case gin.ReleaseMode:
		gin.SetMode(gin.ReleaseMode)
	case gin.TestMode:
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}
	r := gin.New()
	return r
}

// HTTPHandler returns http handler
func (h *Handler) HTTPHandler() (*gin.Engine, error) {
	router := newRouter(h.Mode)

	router.HandleMethodNotAllowed = true

	router.Use(
		gin.Logger(),
		recovery.PanicRecovery(),
		cors.Default(),
		location.Default(),
		requestid.RequestID(),
	)

	// system
	router.GET("/health", h.health)
	router.GET("/readiness", h.readiness)

	// apps

	return router, nil
}

func (h *Handler) health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "ok"})
}

func (h *Handler) readiness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "readiness"})
}
