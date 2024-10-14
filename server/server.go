package server

import (
	"log"
	"net"

	customerpb2 "github.com/nicce/go-grpc-lab/customerpb"
	"github.com/nicce/go-grpc-lab/customers"
	"github.com/nicce/go-grpc-lab/provider/customer"

	"google.golang.org/grpc"
)

// Server struct for the gRPC service.
type Server struct {
	grpcServer *grpc.Server
	listener   net.Listener
}

// Config for the server.
type Config struct {
	Port                   string
	Customers              []*customers.Customer
	MaxDelayInMilliseconds int
}

// New creates a new server.
func New(cfg Config) *Server {
	listener, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	customerProvider := customer.New(cfg.Customers, cfg.MaxDelayInMilliseconds)

	grpcServer := grpc.NewServer()
	customerService := customers.New(customerProvider)
	customerpb2.RegisterCustomerServiceServer(grpcServer, customerService)

	return &Server{
		grpcServer: grpcServer,
		listener:   listener,
	}
}

// Serve starts the server.
func (s *Server) Serve() error {
	return s.grpcServer.Serve(s.listener)
}
