package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/lcolman/fabrikam-auth-poc/internal/config"
	"github.com/rs/cors"
)

type Server interface {
	Start(address string, port int)
	AddRoute(route Route)
	AddMiddleware(middleware mux.MiddlewareFunc)
	Logger() *log.Logger
}

type apiServer struct {
	base        string
	crossOrigin *cors.Cors
	logger      *log.Logger
	router      *mux.Router
}

func NewServer(cfg config.ApplicationConfig, middlewares ...mux.MiddlewareFunc) Server {
	router := mux.NewRouter().PathPrefix(cfg.BasePath()).Subrouter()
	router.Use(middlewares...)
	return apiServer{
		base:        cfg.BasePath(),
		crossOrigin: cors.AllowAll(),
		logger:      log.New(os.Stderr, fmt.Sprintf("[%s] ", strings.ToUpper(cfg.LoggingName())), log.Ltime),
		router:      router,
	}
}

func (s apiServer) Logger() *log.Logger {
	return s.logger
}

func (s apiServer) AddRoute(route Route) {
	s.router.NewRoute().Methods(route.Method).Path(route.Path).HandlerFunc(route.Handler)
	// also register a handler for the cross origin resource sharing pre-flight OPTIONS request:
	s.router.NewRoute().Methods(http.MethodOptions).Path(route.Path).HandlerFunc(s.crossOrigin.HandlerFunc)
}

func (s apiServer) AddMiddleware(middleware mux.MiddlewareFunc) {
	s.router.Use(middleware)
}

func (s apiServer) Start(address string, port int) {
	loggingHost := address
	if address == "" {
		loggingHost = "localhost"
	}
	s.logger.Printf("INFO Starting server on http://%s:%v%s/", loggingHost, port, s.base)
	err := http.ListenAndServe(fmt.Sprintf("%s:%v", address, port), s.router)
	panic(err)
}
