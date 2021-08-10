package checkout

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gitlab.com/lucafmarques/hash-test/auth"
	"gitlab.com/lucafmarques/hash-test/config"
	"gitlab.com/lucafmarques/hash-test/discount"
	_ "gitlab.com/lucafmarques/hash-test/docs"
	"gitlab.com/lucafmarques/hash-test/repository"
)

var ALLOW_DOCS_MODES = map[string]bool{
	"DEVELOPMENT": true,
	"STAGING":     true,
	"PRODUCTION":  false,
}

type Service struct {
	Server echo.Echo
	Core   Core
	Config config.ServiceConfig
}

func NewCheckoutService(config config.ServiceConfig, client discount.DiscountClient, repo repository.Repository) *Service {
	server := echo.New()
	server.HideBanner = true
	server.HidePort = true

	if config.Port == "" {
		config.Port = ":8080"
		log.Info("no checkout port config, defaulting to :8080")
	}

	if config.Environment == "" {
		config.Environment = "DEVELOPMENT"
	}

	return &Service{
		Server: *server,
		Core:   NewCore(config.Core, client, repo),
		Config: config,
	}
}

func (svc *Service) Start() {
	go func() {
		if err := svc.Server.Start(svc.Config.Port); err != nil && err != http.ErrServerClosed {
			svc.Server.Logger.Fatal(err)
		}
	}()

	log.Info("Starting checkout server")
}

func (svc *Service) Stop() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := svc.Server.Shutdown(ctx); err != nil {
		svc.Server.Logger.Fatal(err)
	}

	log.Info("Finishing gracefully shutting down service")
}

func (svc *Service) RegisterRoutes() {
	svc.Server.GET("/products", svc.GetAllProducts, auth.AuthMiddleware())
	svc.Server.GET("/discount/:id", svc.GetProductDiscount, auth.AuthMiddleware())
	svc.Server.POST("/checkout", svc.PostCheckout, auth.AuthMiddleware())

	if ok := ALLOW_DOCS_MODES[svc.Config.Environment]; ok {
		svc.Server.GET("/docs/*", echoSwagger.WrapHandler)
	}
}

func (svc *Service) ApplyMiddlewares(middewares ...echo.MiddlewareFunc) {
	svc.Server.Use(middewares...)
}
