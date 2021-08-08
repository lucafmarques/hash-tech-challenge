package discount

import (
	context "context"
	"fmt"
	"time"

	"gitlab.com/lucafmarques/hash-test/config"
	grpc "google.golang.org/grpc"
)

type Discount struct {
	Client DiscountClient
	Config config.DiscountConfig
}

func NewDiscountConn(config config.DiscountConfig) (*grpc.ClientConn, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(config.Timeout))

	opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()}
	conn, err := grpc.DialContext(ctx, config.Host, opts...)
	if err != nil {
		cancel()
		return nil, nil, fmt.Errorf("fail to dial: %v", err)
	}

	return conn, cancel, nil
}
