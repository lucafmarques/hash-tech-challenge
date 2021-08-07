package main

import (
	"context"
	_ "embed"
	"log"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gitlab.com/lucafmarques/hash-test/auth"
	"gitlab.com/lucafmarques/hash-test/checkout"
	"gitlab.com/lucafmarques/hash-test/discount"
	"gitlab.com/lucafmarques/hash-test/repository"
	"google.golang.org/grpc"
)

var (
	discountServerAddr = "discount:50051"
)

func main() {
	grpcOpts := []grpc.DialOption{
		grpc.WithInsecure(), grpc.WithBlock(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	conn, err := grpc.DialContext(ctx, discountServerAddr, grpcOpts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	server := echo.New()
	client := discount.NewDiscountClient(conn)
	repo, err := repository.NewEmbedRepository()
	if err != nil {
		log.Fatalf("Failed loading service repository: %s", err)
	}

	svc := checkout.NewCheckoutService(*server, client, repo)
	defer svc.Stop()

	svc.ApplyMiddlewares(middleware.Logger(), middleware.Recover(), auth.AuthMiddleware())
	svc.RegisterRoutes()

	svc.Start()
}
