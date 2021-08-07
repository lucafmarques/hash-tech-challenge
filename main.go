package main

import (
	"context"
	_ "embed"
	"log"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gitlab.com/lucafmarques/hash-test/checkout"
	"gitlab.com/lucafmarques/hash-test/discount"
	"gitlab.com/lucafmarques/hash-test/repository"
	"google.golang.org/grpc"
)

var (
	discountServerAddr = "discount:50051"
)

// @title Hash's Checkout API
// @version 1.0.0
// @description API for receiving cart info and generating a checkout order with proper discounts received by calling Discount service.
// @contact.name Luca F. Marques
// @contact.email lucafmarqs@gmail.com
// @license.name MIT

func main() {
	config := &Config{}
	err := config.LoadFromEnv("CONFIG_PATH", "config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	grpcOpts := []grpc.DialOption{
		grpc.WithInsecure(), grpc.WithBlock(),
	}

	conn, err := grpc.DialContext(ctx, discountServerAddr, grpcOpts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	server := echo.New()
	client := discount.NewDiscountClient(conn)
	repo, err := repository.NewEmbedRepository()
	if err != nil {
		log.Fatalf("failed loading service repository: %s", err)
	}

	svc := checkout.NewCheckoutService(*server, client, repo)
	defer svc.Stop()

	svc.ApplyMiddlewares(middleware.Logger(), middleware.Recover())
	svc.RegisterRoutes()

	svc.Start()
}
