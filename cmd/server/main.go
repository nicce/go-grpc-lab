package main

import (
	"context"
	"github.com/nicce/go-grpc-lab/internal/customers"
	"github.com/nicce/go-grpc-lab/internal/environment"
	"github.com/nicce/go-grpc-lab/internal/xgrpc"
	"log"

	"github.com/brianvoe/gofakeit/v7"
)

func main() {
	ctx := context.Background()
	env := environment.GetEnvironment(ctx)

	log.Printf("Starting server with cfg: %+v\n", env)

	s := xgrpc.New(xgrpc.Config{Port: env.Port, Customers: getMockCustomers(env.NumberOfCustomers), MaxDelayInMilliseconds: env.MaxDelayInMilliseconds})
	if err := s.Serve(); err != nil {
		log.Fatal(err)
	}
}

func getMockCustomers(nbrOfCustomers int) []*customers.Customer {
	faker := gofakeit.New(0)
	c := make([]*customers.Customer, 0)

	for i := 0; i < nbrOfCustomers; i++ {
		customer := &customers.Customer{
			ID:    faker.UUID(),
			Name:  faker.Name(),
			Email: faker.Email(),
			Phone: faker.Phone(),
			Address: customers.Address{
				Street: faker.Street(),
				City:   faker.City(),
				State:  faker.State(),
				Zip:    faker.Zip(),
			},
		}
		c = append(c, customer)
	}

	return c
}
