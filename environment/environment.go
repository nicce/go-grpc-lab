package environment

import (
	"context"
	"log"

	"github.com/nicce/go-grpc-lab/xenvironment"
)

const (
	defaultMaxDelaysInMilliseconds = 500
	defaultNumberOfCustomers       = 10
	defaultPort                    = "8080"
)

// Environment - the struct containing the environment information.
type Environment struct {
	MaxDelayInMilliseconds int    `env:"MAX_DELAY_IN_MILLISECONDS"`
	NumberOfCustomers      int    `env:"NUMBER_OF_CUSTOMERS"`
	Port                   string `env:"PORT"`
}

// GetEnvironment - returns the application environment.
func GetEnvironment(_ context.Context) *Environment {
	env := &Environment{
		MaxDelayInMilliseconds: defaultMaxDelaysInMilliseconds,
		NumberOfCustomers:      defaultNumberOfCustomers,
		Port:                   defaultPort,
	}

	if err := xenvironment.GetEnvironment(env); err != nil {
		log.Fatalf("Error reading the environment variables: %v", err)
	}

	return env
}
