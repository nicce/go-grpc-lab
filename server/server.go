package server

import (
	"github.com/nicce/go-grpc-lab/api/customerpb"
	"github.com/nicce/go-grpc-lab/api/customers"
	"log"
	"net"

	"google.golang.org/grpc"
)

// Server struct for the gRPC service.
type Server struct {
	grpcServer *grpc.Server
	listener   net.Listener
}

// Config for the server.
type Config struct {
	Port      string
	Customers []*customerpb.GetCustomerResponse
}

// New creates a new server.
func New(cfg Config) *Server {
	listener, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	grpcServer := grpc.NewServer()
	customerService := customers.New(cfg.Customers)
	customerpb.RegisterCustomerServiceServer(grpcServer, customerService)

	return &Server{
		grpcServer: grpcServer,
		listener:   listener,
	}
}

// Serve starts the server.
func (s *Server) Serve() error {
	return s.grpcServer.Serve(s.listener)
}
