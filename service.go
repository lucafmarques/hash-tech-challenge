package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"gitlab.com/lucafmarques/hash-test/discount"
	"gitlab.com/lucafmarques/hash-test/repository"
)

type Service struct {
	Server         echo.Echo
	DiscountClient discount.DiscountClient
	Repository     repository.Repository
}

func NewCheckoutService(server echo.Echo, client discount.DiscountClient, repo repository.Repository) *Service {
	server.Logger = log.New("checkout")

	server.HideBanner = true
	server.HidePort = true

	server.Logger.Info("Starting checkout server")

	return &Service{
		Server:         server,
		DiscountClient: client,
		Repository:     repo,
	}
}

func (svc *Service) Start() {
	go func() {
		if err := svc.Server.Start(":8080"); err != nil && err != http.ErrServerClosed {
			svc.Server.Logger.Fatal(err)
		}
	}()
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
}

func (svc *Service) RegisterRoutes() {
	svc.Server.GET("/hello", svc.HandleHello)
	svc.Server.GET("/products", svc.GetAllProducts)
	svc.Server.GET("/discount/:id", svc.GetProductDiscount)
	svc.Server.POST("/checkout", svc.PostCheckout)
}

func (svc *Service) ApplyMiddlewares(middewares ...echo.MiddlewareFunc) {
	svc.Server.Use(middewares...)
}
