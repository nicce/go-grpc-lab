package customer

import (
	"context"
	"math/rand"
	"time"

	"github.com/nicce/go-grpc-lab/internal/customers"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MockProvider describes the MockProvider struct.
type MockProvider struct {
	customers []*customers.Customer
	maxDelay  int
}

// New creates a new mock provider.
func New(customers []*customers.Customer, maxDelay int) *MockProvider {
	return &MockProvider{
		customers: customers,
		maxDelay:  maxDelay,
	}
}

// GetCustomer returns a customer by ID.
func (p *MockProvider) GetCustomer(_ context.Context, id string) (*customers.Customer, error) {
	for _, c := range p.customers {
		if c.ID == id {
			return c, nil
		}
	}

	return nil, status.Errorf(codes.NotFound, "customer not found")
}

// StreamCustomers returns a channel of customers.
func (p *MockProvider) StreamCustomers(_ context.Context, ids []string) (<-chan *customers.Customer, error) {
	customerChan := make(chan *customers.Customer)

	go func() {
		for _, c := range p.customers {
			if shouldSendCustomer(ids, c) {
				customerChan <- c

				// delay for simulation purposes
				time.Sleep(time.Duration(rand.Intn(p.maxDelay)) * time.Millisecond) //nolint:gosec
			}
		}

		// Close channel once all customers has been processed
		close(customerChan)
	}()

	return customerChan, nil
}

func shouldSendCustomer(ids []string, c *customers.Customer) bool {
	if len(ids) == 0 {
		return true
	}

	for _, id := range ids {
		if c.ID == id {
			return true
		}
	}

	return false
}
