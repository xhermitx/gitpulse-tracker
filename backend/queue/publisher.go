package queue

type Publisher interface {
	Publish(data any, queueName string) error
}
