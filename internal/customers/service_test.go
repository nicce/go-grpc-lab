package customers_test

import (
	"context"
	"github.com/nicce/go-grpc-lab/api/gen/customerpb"
	"github.com/nicce/go-grpc-lab/internal/customers"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

type mockProvider struct {
	// should have a list of mock data provided
}

func (m *mockProvider) GetCustomer(ctx context.Context, id string) (*customers.Customer, error) {
	return &customers.Customer{
		ID:      "1",
		Name:    "John Doe",
		Email:   "",
		Phone:   "",
		Address: customers.Address{},
	}, nil
}

func (m *mockProvider) StreamCustomers(ctx context.Context, ids []string) (<-chan *customers.Customer, error) {
	customerChan := make(chan *customers.Customer)
	defer close(customerChan)

	customerChan <- &customers.Customer{
		ID:      "1",
		Name:    "John Doe",
		Email:   "",
		Phone:   "",
		Address: customers.Address{},
	}

	return customerChan, nil
}

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func startTestGRPCServer(provider *mockProvider) *grpc.Server {
	s := grpc.NewServer()
	// Register the service with your gRPC server
	customerGRPCService := customers.New(provider)
	customerpb.RegisterCustomerServiceServer(s, customerGRPCService)

	go func() {
		if err := s.Serve(lis); err != nil {
			panic("Server failed to start: " + err.Error())
		}
	}()
	return s
}

func TestGetCustomer_Success(t *testing.T) {
	// Create a mock provider
	provider := &mockProvider{}

	// Start the gRPC server with the mock provider
	server := startTestGRPCServer(provider)
	defer server.Stop()

	// Create a connection using bufcon
	ctx := context.Background()
	conn, err := grpc.NewClient("passthrough://bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(conn)

	assert.NoError(t, err)

	// Create the client
	client := customerpb.NewCustomerServiceClient(conn)

	// Call ListCustomers via gRPC client
	c, err := client.GetCustomer(ctx, &customerpb.GetCustomerRequest{Id: "1"})
	assert.NoError(t, err)

	// Validate the received response
	assert.Equal(t, "1", c.Id)
	assert.Equal(t, "John Doe", c.Name)
}
