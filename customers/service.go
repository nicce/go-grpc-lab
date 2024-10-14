package customers

import (
	"context"

	"github.com/nicce/go-grpc-lab/customerpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Service a customer service.
type Service struct {
	provider provider
	customerpb.UnimplementedCustomerServiceServer
}

// Customer describes the Customer struct.
type Customer struct {
	ID      string
	Name    string
	Email   string
	Phone   string
	Address Address
}

// Address describes the Address struct.
type Address struct {
	Street string
	City   string
	State  string
	Zip    string
}

type provider interface {
	GetCustomer(ctx context.Context, id string) (*Customer, error)
	StreamCustomers(ctx context.Context, ids []string) (<-chan *Customer, error)
}

// New creates a new customer service.
func New(provider provider) *Service {
	return &Service{
		provider: provider,
	}
}

// GetCustomer returns a customer by ID.
func (s *Service) GetCustomer(ctx context.Context, req *customerpb.GetCustomerRequest) (*customerpb.GetCustomerResponse, error) {
	customer, err := s.provider.GetCustomer(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	address := customerpb.AddressResponse{
		Street: customer.Address.Street,
		City:   customer.Address.City,
		State:  customer.Address.State,
		Zip:    customer.Address.Zip,
	}

	return &customerpb.GetCustomerResponse{
		Id:      customer.ID,
		Name:    customer.Name,
		Email:   customer.Email,
		Phone:   customer.Phone,
		Address: &address,
	}, nil
}

// ListCustomers returns a list of customers.
func (s *Service) ListCustomers(req *customerpb.ListCustomersRequest, stream grpc.ServerStreamingServer[customerpb.GetCustomerResponse]) error {
	customerChan, err := s.provider.StreamCustomers(context.Background(), req.Ids)
	if err != nil {
		return err
	}

	for {
		select {
		case customer, ok := <-customerChan:
			if !ok { // Channel closed
				return nil
			}

			address := customerpb.AddressResponse{
				Street: customer.Address.Street,
				City:   customer.Address.City,
				State:  customer.Address.State,
				Zip:    customer.Address.Zip,
			}
			c := customerpb.GetCustomerResponse{
				Id:      customer.ID,
				Name:    customer.Name,
				Email:   customer.Email,
				Phone:   customer.Phone,
				Address: &address,
			}

			if err := stream.Send(&c); err != nil {
				return status.Errorf(codes.Internal, "error sending customer: %v", err)
			}
		case <-stream.Context().Done():
			// Handle client disconnects
			return stream.Context().Err()
		}
	}
}
