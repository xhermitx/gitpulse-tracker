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

func CreateChannel(ch *amqp.Channel) <-chan amqp.Delivery {
	q, err := ch.QueueDeclare(
		models.GITHUB_DATA_QUEUE, // name
		false,                    // durable
		false,                    // delete when unused
		false,                    // exclusive
		false,                    // no-wait
		nil,                      // arguments
	)
	LogError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	LogError(err, "Failed to register a consumer")
	return msgs
}

func HandleQueueData(msgs <-chan amqp.Delivery) {

	if msgs == nil {
		log.Println("empty body from queue")
	}

	for d := range msgs {
		var data models.Candidate
		err := json.Unmarshal(d.Body, &data)
		LogError(err, "failed to read data from queue")

		rdb := redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS"),
			Password: "",
			DB:       0,
		})

		ctx := context.Background()
		rdbClient := redisdb.NewRedisClient(rdb)

		if !data.Status {
			err := rdbClient.Set(ctx, data.TopCandidates)
			LogError(err, "Failed to store the candidate on Redis")

		} else if data.Status {

			fmt.Println("End sequence initiated for : ", data.TopCandidates.JobId)

			// RETRIEVE THE TOP 5 CANDIDATES
			topCandidates, err := rdbClient.Get(ctx, data.TopCandidates.JobId)
			LogError(err, "Failed to Retrieve data from Redis")

			// PUSH TO DB
			err = store.InsertData(topCandidates)
			if err != nil {
				LogError(err, err.Error())
			} else {
				// UPDATE THE STATUS THAT PROFILING IS COMPLETE FOR THE CORRESPONDING JOB ID
				conn, err := rabbitmq.Connect()
				if err != nil {
					LogError(err, "faild to connect to rabbitmq")
				}
				rmqClient := rabbitmq.NewRabbitClient(conn)
				status := models.StatusQueue{
					JobId:  data.TopCandidates.JobId,
					Status: false,
				}
				err = rmqClient.Publish(status, models.STATUS_QUEUE)
				LogError(err, "failed to update profiling status on queue")
			}
		}
	}
}

func LogError(err error, msg string) {
	if err != nil {
		log.Printf("\n%s: %s", msg, err)
	}
}
