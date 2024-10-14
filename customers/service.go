package customers

import (
	"context"
	"github.com/nicce/go-grpc-lab/api/customerpb"
	"google.golang.org/grpc"
)

// Service a customer service.
type Service struct {
	provider CustomerProvider
	customerpb.UnimplementedCustomerServiceServer
}

type Customer struct {
	Id      string
	Name    string
	Email   string
	Phone   string
	Address Address
}

type Address struct {
	Street string
	City   string
	State  string
	Zip    string
}

type CustomerProvider interface {
	GetCustomer(ctx context.Context, id string) (Customer, error)
	ListCustomers(ctx context.Context, ids []string) chan Customer
}

// New creates a new customer service.
func New(provider CustomerProvider) *Service {
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
		Id:      customer.Id,
		Name:    customer.Name,
		Email:   customer.Email,
		Phone:   customer.Phone,
		Address: &address,
	}, nil
}

// ListCustomers returns a list of customers.
func (s *Service) ListCustomersX(req *customerpb.ListCustomersRequest, stream grpc.ServerStreamingServer[customerpb.GetCustomerResponse]) error {
	customers := s.provider.ListCustomers(context.Background(), req.Ids)

	return nil
}
