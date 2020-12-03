package main

import (
	"log"

	"github.com/lcolman/fabrikam-auth-poc/internal/auth"
	"github.com/lcolman/fabrikam-auth-poc/internal/config"
	"github.com/lcolman/fabrikam-auth-poc/internal/screening"
	"github.com/lcolman/fabrikam-auth-poc/internal/server"
)

type serviceProvider struct {
	screeningService screening.Service
	logger           *log.Logger
}

func (sp *serviceProvider) ScreeningService() screening.Service {
	return sp.screeningService
}

func (sp *serviceProvider) Logger() *log.Logger {
	return sp.logger
}

func main() {
	conf := config.NewApplicationConfig()
	server := server.NewServer(conf)
	server.AddMiddleware(auth.Middleware(conf, server.Logger()))
	screeningDB := screening.NewStorageService()
	screeningService := screening.NewScreeningService(screeningDB)
	collaborators := &serviceProvider{screeningService, server.Logger()}
	server.AddRoute(screening.GetRoute(collaborators))
	server.AddRoute(screening.PostRoute(collaborators))
	server.Start("", conf.ServerPort())
}
