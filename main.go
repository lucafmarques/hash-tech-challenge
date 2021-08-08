package main

import (
	_ "embed"
	"log"

	"github.com/labstack/echo/v4/middleware"
	"gitlab.com/lucafmarques/hash-test/checkout"
	"gitlab.com/lucafmarques/hash-test/config"
	"gitlab.com/lucafmarques/hash-test/discount"
	"gitlab.com/lucafmarques/hash-test/repository"
)

// @title Hash's Checkout API
// @version 1.0.0
// @description API for receiving cart info and generating a checkout order with proper discounts received by calling Discount service.
// @contact.name Luca F. Marques
// @contact.email lucafmarqs@gmail.com
// @license.name MIT

func main() {
	config := &config.Config{}
	err := config.LoadFromEnv("CONFIG_PATH", "config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	conn, cancel, err := discount.NewDiscountConn(config.Discount)
	if err != nil {
		log.Fatalf("failed creating discount grpc conn: %v", err)
	}
	defer conn.Close()
	defer cancel()

	repo, err := repository.NewEmbedRepository(config.Repository)
	if err != nil {
		log.Fatalf("failed loading service repository: %s", err)
	}

	client := discount.NewDiscountClient(conn)
	svc := checkout.NewCheckoutService(config.Service, client, repo)
	defer svc.Stop()

	svc.ApplyMiddlewares(middleware.Logger(), middleware.Recover())
	svc.RegisterRoutes()

	svc.Start()
}
