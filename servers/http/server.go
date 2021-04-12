package httpserver

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
)

var (
	// ErrHTTPHandlerDoesNotSet denotes error without setting HTTP handler
	ErrHTTPHandlerDoesNotSet error = errors.New("http(s) handler does not set")
)

type sig string

func (s sig) String() string {
	return "run server error, because: " + string(s)
}

func (sig) Signal() {
	// this is blank
}

var sigRunServer sig

// Option ...
type Option func(*Server) error

// WithHTTPS ...
func WithHTTPS(addr, certPath, keyPath string) Option {
	return func(s *Server) error {
		s.HTTPS.Addr = addr
		s.HTTPS.CertPath = certPath
		s.HTTPS.KeyPath = keyPath
		return nil
	}
}

// Server handles HTTP
type Server struct {
	Handler  http.Handler
	HTTPAddr string
	HTTPS    struct {
		Addr, CertPath, KeyPath string
	}
}

// New returns http server
func New(httpAddr string, handler http.Handler, opts ...Option) (*Server, error) {
	if handler == nil {
		return nil, ErrHTTPHandlerDoesNotSet
	}

	s := &Server{HTTPAddr: httpAddr, Handler: handler}

	for _, opt := range opts {
		opt(s)
	}
	return s, nil
}

// Run runs HTTP and HTTPs servers
func (s *Server) Run() error {
	srv := &http.Server{
		Addr:    s.HTTPAddr,
		Handler: s.Handler,
		//IdleTimeout:       30 * time.Second,
		// ReadTimeout:       10 * time.Second,

		//ReadHeaderTimeout: 250 * time.Millisecond,
		//WriteTimeout:      500 * time.Millisecond,
	}

	quit := make(chan os.Signal, 1)

	// HTTP
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			sigRunServer = sig("ListenAndServe(HTTP):" + err.Error())
			quit <- sigRunServer
		}
	}()

	// HTTPS; TODO move out
	var srvTLS *http.Server = nil
	if s.checkIfSetHTTPS() {
		srvTLS = &http.Server{
			Addr:    s.HTTPS.Addr,
			Handler: s.Handler,
		}
		go func() {
			if err := srvTLS.ListenAndServeTLS(s.HTTPS.CertPath, s.HTTPS.KeyPath); err != nil {
				sigRunServer = sig("ListenAndServeTLS(HTTPs):" + err.Error())
				quit <- sigRunServer
			}
		}()
	}

	// graceful shutdown
	signal.Notify(quit, os.Kill, syscall.SIGINT, syscall.SIGTERM, sigRunServer)
	v := <-quit
	var runSrv error
	if s, ok := v.(sig); ok {
		runSrv = errors.New(s.String())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	if srvTLS != nil {
		srvTLS.Shutdown(ctx)
	}

	if runSrv != nil {
		return runSrv
	}
	return nil

}

func (s *Server) checkIfSetHTTPS() bool {
	if len(s.HTTPS.Addr) > 0 && len(s.HTTPS.CertPath) > 0 && len(s.HTTPS.KeyPath) > 0 {
		return true
	}
	return false
}
