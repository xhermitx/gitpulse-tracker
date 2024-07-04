package servers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
	redis "github.com/redis/go-redis/v9"
	"github.com/xhermitx/gitpulse-tracker/profiler-service/api/rabbitmq"
	"github.com/xhermitx/gitpulse-tracker/profiler-service/api/redisdb"
	"github.com/xhermitx/gitpulse-tracker/profiler-service/models"
	"github.com/xhermitx/gitpulse-tracker/profiler-service/store"
)

func Listener() {
	fmt.Println("RABBITMQ: ", os.Getenv("RABBITMQ"))

	conn, err := amqp.Dial(os.Getenv("RABBITMQ"))
	logError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	logError(err, "Failed to open a channel")
	defer ch.Close()

	msgs := createChannel(ch)

	forever := make(chan struct{})

	go func() {
		for d := range msgs {
			var data models.Candidate
			err := json.Unmarshal(d.Body, &data)
			logError(err, err.Error())

			rdb := redis.NewClient(&redis.Options{
				Addr:     os.Getenv("REDIS"),
				Password: "",
				DB:       0,
			})

			ctx := context.Background()

			client := redisdb.NewRedisClient(rdb)

			if !data.Status {
				err := client.Set(ctx, data.TopCandidates)
				logError(err, "Failed to store the candidate on Redis")

			} else if data.Status {

				fmt.Println("End sequence initiated for : ", data.TopCandidates.JobId)

				// RETRIEVE THE TOP 5 CANDIDATES
				topCandidates, err := client.Get(ctx, data.TopCandidates.JobId)
				logError(err, "Failed to Retrieve data from Redis")

				// PUSH TO DB
				err = store.InsertData(topCandidates)
				logError(err, err.Error())

				if err == nil {
					conn, err2 := rabbitmq.Connect()
					logError(err2, err2.Error())

					redisClient := rabbitmq.NewRedisClient(conn)

					status := models.StatusQueue{
						JobId:  data.TopCandidates.JobId,
						Status: false,
					}

					err3 := redisClient.Publish(status, models.STATUS_QUEUE)
					logError(err3, err3.Error())
				}
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func createChannel(ch *amqp.Channel) <-chan amqp.Delivery {

	q, err := ch.QueueDeclare(
		models.GITHUB_DATA_QUEUE, // name
		false,                    // durable
		false,                    // delete when unused
		false,                    // exclusive
		false,                    // no-wait
		nil,                      // arguments
	)
	logError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	logError(err, "Failed to register a consumer")

	return msgs
}

func logError(err error, msg string) {
	if err != nil {
		log.Printf("\n%s: %s", msg, err)
	}
}
