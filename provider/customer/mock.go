package customer

import (
	"context"
	"github.com/nicce/go-grpc-lab/customers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type CustomerProvider struct {
	customers []*customers.Customer
}

func New(customers []*customers.Customer) *CustomerProvider {
	return &CustomerProvider{
		customers: customers,
	}
}

func (p *CustomerProvider) GetCustomer(_ context.Context, id string) (*customers.Customer, error) {
	for _, c := range p.customers {
		if c.Id == id {
			return c, nil
		}
	}

	return nil, status.Errorf(codes.NotFound, "customer not found") // TODO: proper error handling that later translates to a gRPC error
}

func (p *CustomerProvider) StreamCustomers(_ context.Context, ids []string) (<-chan *customers.Customer, error) {
	customerChan := make(chan *customers.Customer)

	go func() {
		for _, c := range p.customers {
			if shouldSendCustomer(ids, c) {
				customerChan <- c
				time.Sleep(500 * time.Millisecond) // delay for simulation purposes
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
		if c.Id == id {
			return true
		}
	}

	return false
}
