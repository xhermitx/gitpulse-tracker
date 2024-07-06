package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitClient struct {
	conn *amqp.Connection
}

func NewRabbitClient(connection *amqp.Connection) *RabbitClient {
	return &RabbitClient{
		conn: connection,
	}
}

func (client *RabbitClient) Publish(data any, queueName string) error {

	ch, err := client.conn.Channel()
	if err != nil {
		log.Println("failed to create channel")
		return err
	}

	q, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("failed to declare queue")
		return err
	}

	var timer time.Time

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	ctx = context.WithValue(ctx, timer, time.Now())
	defer cancel()

	body, err := json.Marshal(data)
	if err != nil {
		log.Println("failed to marshal data")
		return err
	}

	if err = ch.PublishWithContext(ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		}); err != nil {
		log.Println("failed to publish data")
		return err
	}

	log.Println("DATA PUBLISHED !")
	return nil
}

func Connect() (*amqp.Connection, error) {
	var (
		counts     int64
		backOff    = 1 * time.Second
		connection *amqp.Connection
	)

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
			log.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(2, float64(counts))) * time.Second
		log.Println("backing off...")
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
