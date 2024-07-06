package servers

import (
	"fmt"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Listener() {
	fmt.Println("RABBITMQ: ", os.Getenv("RABBITMQ"))

	conn, err := amqp.Dial(os.Getenv("RABBITMQ"))
	LogError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	LogError(err, "Failed to open a channel")
	defer ch.Close()

	msgs := CreateChannel(ch)

	go HandleQueueData(msgs)

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	select {}
}
