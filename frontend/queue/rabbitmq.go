package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/xhermitx/gitpulse-tracker/frontend/models"
)

type RabbitMQ struct {
	Data      any
	QueueName string
}

func NewRabbitMQClient(data any, queueName string) *RabbitMQ {
	return &RabbitMQ{
		Data:      data,
		QueueName: queueName,
	}
}

// FUNCTION TO PUSH CANDIDATE DATA TO THE MESSAGE QUEUE
func (mq *RabbitMQ) Publish() error {

	conn, err := connect()
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		models.STATUS_QUEUE, // name
		false,               // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)

	failOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(mq.Data)

	failOnError(err, fmt.Sprintf("Failed to Parse Status for: %d", mq.Data.(models.StatusQueue).JobId))

	if err = ch.PublishWithContext(ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		}); err != nil {
		return err
	}

	log.Printf("\n[x] Sent status for job %d", mq.Data.(models.StatusQueue).JobId)

	return nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// RETRY CONNECTION WITH EXPONENTIAL TIMEOUT
func connect() (*amqp.Connection, error) {
	var (
		counts     int64
		backOff    = 1 * time.Second
		connection *amqp.Connection
	)

	log.Println(os.Getenv("RABBITMQ"))

	for {
		c, err := amqp.Dial(os.Getenv("RABBITMQ"))
		if err != nil {
			fmt.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			log.Println("Connected to RabbitMQ!")
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(2, float64(counts))) * time.Second
		log.Println("backing off...")
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
