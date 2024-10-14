package main

import (
	"github.com/nicce/go-grpc-lab/api/customerpb"
	"github.com/nicce/go-grpc-lab/api/server"
	"log"
	"os"

	"github.com/brianvoe/gofakeit/v7"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	s := server.New(server.Config{Port: port, Customers: getMockCustomers()})
	if err := s.Serve(); err != nil {
		log.Fatal(err)
	}
}

func getMockCustomers() []*customerpb.GetCustomerResponse {
	faker := gofakeit.New(0)
	customers := make([]*customerpb.GetCustomerResponse, 0)

	for i := 0; i < 100; i++ {
		customer := &customerpb.GetCustomerResponse{
			Id:    faker.UUID(),
			Name:  faker.Name(),
			Email: faker.Email(),
			Phone: faker.Phone(),
			Address: &customerpb.AddressResponse{
				Street: faker.Street(),
				City:   faker.City(),
				State:  faker.State(),
				Zip:    faker.Zip(),
			},
		}
		customers = append(customers, customer)
	}

	return customers
}
