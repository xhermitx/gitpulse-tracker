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

type RedisClient struct {
	conn *amqp.Connection
}

func NewRedisClient(connection *amqp.Connection) *RedisClient {
	return &RedisClient{
		conn: connection,
	}
}

func (client *RedisClient) Publish(data any, queueName string) error {

	ch, err := client.conn.Channel()
	failOnError(err, "failed to open a channel")

	q, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "failed to declare the queue")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(data)
	failOnError(err, err.Error())

	err = ch.PublishWithContext(ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		})
	failOnError(err, "failed to publish data")

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

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Printf("\n%s: %s", msg, err)
	}
}
