package httphandler

import (
	"context"
	"net/http"
	"simpleBackend/handlers/http/middleware/cors"
	"simpleBackend/handlers/http/middleware/recovery"
	"simpleBackend/handlers/http/middleware/requestid"
	"simpleBackend/handlers/maindb"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/pkg/errors"

	"github.com/gin-contrib/location"
	// https://github.com/gin-contrib/pprof
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Option is option for handler
type Option func(*Handler) error

var (
	ErrEmptyNasaAPIKey      error = errors.New("Nasn API key is empty")
	ErrMainDBDoesNotExist   error = errors.New("main database does not exist")
	ErrEmptyRedisServerAddr error = errors.New("address for redis cache server is empty")
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

// WithRedisCache sets nasa api key
func WithRedisCache(addr, password string) Option {
	return func(h *Handler) error {

		if len(addr) == 0 {
			return ErrEmptyRedisServerAddr
		}

		h.RedisCache = redis.NewClient(&redis.Options{Addr: addr, Password: password})

		timeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		if _, err := h.RedisCache.Ping(timeout).Result(); err != nil {
			return errors.Wrap(err, "failed to ping the redis server")
		}

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

	RedisCache *redis.Client
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

	if h.Mode == "debug" {
		pprof.Register(router)
	}

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

// Destroy destroy the resources is used by this object.
func (h *Handler) Destroy() error {
	// close database connection
	sqlDB, err := h.MainDB.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	defer h.RedisCache.Close()

	return nil
}

func (h *Handler) health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "ok"})
}

func (h *Handler) readiness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "readiness"})
}
