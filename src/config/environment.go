package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type EnvironmentSpec struct {
	// Connection strings
	DbConnectionString       string `split_words:"true" default:"postgres://postgres:postgres@localhost:5432/codelabs?sslmode=disable"`
	RabbitMQConnectionString string `split_words:"true" default:"amqp://rabbitmq:rabbitmq@localhost:5672/"`
}

var environment *EnvironmentSpec

func GetEnvironment() *EnvironmentSpec {
	if environment == nil {
		environment = &EnvironmentSpec{}

		err := envconfig.Process("", environment)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	return environment
}
