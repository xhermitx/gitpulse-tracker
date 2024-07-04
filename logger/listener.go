package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	STATUS_QUEUE   = "profiling_status_queue"
	USERNAME_QUEUE = "username_queue"
)

type StatusQueue struct {
	JobId  uint
	Status bool
}

func Listener() {
	fmt.Println("RABBITMQ: ", os.Getenv("RABBITMQ"))

	conn, err := amqp.Dial(os.Getenv("RABBITMQ"))
	logError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	logError(err, "Failed to open a channel")
	defer ch.Close()

	msgs := createChannel(ch)
	logError(err, "error creating a channel")

	forever := make(chan struct{})

	go func() {
		for d := range msgs {
			var data StatusQueue
			err := json.Unmarshal(d.Body, &data)
			logError(err, "failed to unmarshal queue data")

			var timer time.Time

			if data.Status {
				timer = time.Now()
				_ = timer // To handle the warning. Timer is used for tracking initial and final queue status
			} else {
				totalTime := time.Since(timer)

				f, err := os.OpenFile("logs.md", os.O_CREATE|os.O_APPEND, 0777)
				logError(err, "failed to open the logs")

				content := fmt.Sprintf("Total time taken for the synchronous handling of Drive Data for jobId: %d : %f", data.JobId, totalTime.Seconds())
				f.WriteString(content)
				f.Close()
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func createChannel(ch *amqp.Channel) <-chan amqp.Delivery {

	q, err := ch.QueueDeclare(
		STATUS_QUEUE, // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	logError(err, "failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	logError(err, "failed to register a consumer")

	return msgs
}

func logError(err error, msg string) {
	if err != nil {
		log.Printf("\n%s: %s", msg, err)
	}
}
