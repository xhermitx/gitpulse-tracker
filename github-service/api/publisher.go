package api

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

type Client struct {
	conn *amqp.Connection
}

func NewQueueConnection(conn *amqp.Connection) *Client {
	return &Client{
		conn: conn,
	}
}

// FUNCTION TO PUSH CANDIDATE DATA TO THE MESSAGE QUEUE
func (cq Client) Publish(queueName string, data any) error {

	ch, err := cq.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// DECLARE THE QUEUE
	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(data)
	if err != nil {
		log.Println("failed to marshal data")
		return err
	}

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		})
	if err != nil {
		log.Println("failed to publish data to queue")
		return err
	}

	log.Printf(" [x] Sent %s\n", body)

	return nil
}

// RETRY CONNECTION WITH EXPONENTIAL TIMEOUT
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
