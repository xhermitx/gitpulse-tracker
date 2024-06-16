package servers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
	redis "github.com/redis/go-redis/v9"
	api "github.com/xhermitx/gitpulse-tracker/profiler-service/API"
	"github.com/xhermitx/gitpulse-tracker/profiler-service/models"
	"github.com/xhermitx/gitpulse-tracker/profiler-service/store"
)

func Listener() {
	fmt.Println("RABBITMQ: ", os.Getenv("RABBITMQ"))

	conn, err := amqp.Dial(os.Getenv("RABBITMQ"))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// DECLARE THE QUEUE
	q, err := ch.QueueDeclare(
		"github_data_queue", // name
		false,               // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan struct{})

	go func() {
		for d := range msgs {
			var data models.Candidate
			if err := json.Unmarshal(d.Body, &data); err != nil {
				log.Println(err)
			}

			rdb := redis.NewClient(&redis.Options{
				Addr:     os.Getenv("REDIS"),
				Password: "", // no password set
				DB:       0,  // use default DB
			})

			ctx := context.Background()

			if !data.Status {
				// STORE IN REDIS
				if err := api.Set(data.RedisCandidate, rdb, ctx); err != nil {
					failOnError(err, "Failed to store the candidate on Redis")
				}

			} else if data.Status {

				fmt.Println("End sequence initiated for : ", data.JobId)

				// RETRIEVE THE TOP 5 CANDIDATES
				topCandidates, err := api.Get(data.JobId, rdb, ctx)
				if err != nil {
					failOnError(err, "Failed to Retrieve data from Redis")
				}

				// PUSH TO DB
				if err = store.InsertData(topCandidates); err != nil {
					log.Fatal(err)
				}
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
