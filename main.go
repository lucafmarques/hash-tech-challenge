package main

import (
	_ "embed"
	"log"

	"github.com/labstack/echo/v4/middleware"
	"gitlab.com/lucafmarques/hash-test/auth"
	"google.golang.org/grpc"
)

var (
	//go:embed products.json
	Products           []byte
	discountServerAddr = "localhost:50051"
)

func main() {
	grpcOpts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	conn, err := grpc.Dial(discountServerAddr, grpcOpts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	svc := NewCheckoutService(conn)
	defer svc.Stop()

	svc.ApplyMiddlewares(middleware.Logger(), middleware.Recover(), auth.AuthMiddleware())
	svc.RegisterRoutes()

	svc.Start()
}
