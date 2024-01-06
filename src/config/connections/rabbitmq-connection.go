package connections

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/upb-code-labs/submissions-status-updater-microservice/src/config"
)

var rabbitMQChannel *amqp.Channel

func ConnectToRabbitMQ() {
	// Stablish connection
	rabbitMQConnectionString := config.GetEnvironment().RabbitMQConnectionString
	conn, err := amqp.Dial(rabbitMQConnectionString)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Get channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err.Error())
	}

	// Set channel
	log.Println("[RabbitMQ]: Connected")
	rabbitMQChannel = ch
}

func CloseRabbitMQConnection() {
	if rabbitMQChannel != nil {
		rabbitMQChannel.Close()
	}
}

func GetRabbitMQChannel() *amqp.Channel {
	if rabbitMQChannel == nil {
		ConnectToRabbitMQ()
	}

	return rabbitMQChannel
}
