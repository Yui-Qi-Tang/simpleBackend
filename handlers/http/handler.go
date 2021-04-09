package httphandler

import (
	"net/http"
	"simpleBackend/handlers/http/middleware/cors"
	"simpleBackend/handlers/http/middleware/recovery"
	"simpleBackend/handlers/http/middleware/requestid"
	"simpleBackend/handlers/maindb"
	"time"

	"github.com/pkg/errors"

	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Option is option for handler
type Option func(*Handler) error

var (
	ErrEmptyNasaAPIKey    error = errors.New("Nasn API key is empty")
	ErrMainDBDoesNotExist error = errors.New("main database does not exist")
)

// WithNasaAPIKey sets nasa api key
func WithNasaAPIKey(key string) Option {
	return func(h *Handler) error {

		if len(key) == 0 {
			return ErrEmptyNasaAPIKey
		}
		h.NasaAPIKey = key
		return nil
	}
}

// WithMainDatabase set default database config
func WithMainDatabase(dbType, dsn string, maxOpenConns, maxIdleConns int, connMaxLife time.Duration) Option {
	return func(h *Handler) error {
		db, err := maindb.New(dbType, dsn, maxOpenConns, maxIdleConns, connMaxLife)
		if err != nil {
			return errors.Wrap(err, "failed to config main database")
		}

		h.MainDB = db

		return nil
	}
}

// Handler handles the http service
type Handler struct {
	Mode       string
	NasaAPIKey string

	MainDB *gorm.DB
}

// New returns http handler
func New(mode string, opts ...Option) (*Handler, error) {
	h := &Handler{Mode: mode}

	for _, opt := range opts {
		if err := opt(h); err != nil {
			return nil, err
		}
	}

	if h.MainDB == nil {
		return nil, ErrMainDBDoesNotExist
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
	router.GET("/nasa/apod", h.Apod)

	return router, nil
}

func (h *Handler) health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "ok"})
}

func (h *Handler) readiness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "readiness"})
}
