package discount

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/lucafmarques/hash-test/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const (
	host    = "bufnet"
	bufSize = 1024 * 1024
)

var list *bufconn.Listener

type MockDiscountService struct{}

func (d MockDiscountService) GetDiscount(context.Context, *GetDiscountRequest) (*GetDiscountResponse, error) {
	return &GetDiscountResponse{}, nil
}

func (d MockDiscountService) mustEmbedUnimplementedDiscountServer() {}

func init() {
	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer()
	RegisterDiscountServer(s, &MockDiscountService{})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server existed with error: %v", err)
		}
	}()
}

func dialer(context.Context, string) (net.Conn, error) {
	return list.Dial()
}

func TestNewDiscountConn(t *testing.T) {
	config := config.DiscountConfig{
		Host:    host,
		Timeout: 10,
	}

	opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithContextDialer(dialer)}

	conn, _, err := NewDiscountConn(config, opts)
	assert.Nil(t, err, "Failed asserting non-error when creating connection")
	assert.IsType(t, &grpc.ClientConn{}, conn, "Failed asserting conn type")
}

func TestNewDiscountConnError(t *testing.T) {
	config := config.DiscountConfig{
		Host:    host,
		Timeout: 10,
	}

	conn, _, err := NewDiscountConn(config, nil)
	assert.Error(t, err, "Failed asserting error when creating connection")
	assert.Nil(t, conn, "Failed asserting conn type")
}
