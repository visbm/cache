package rabbitmq

import (
	"cache/internal/config"
	"cache/pkg/logger"
	"fmt"
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConsumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   *amqp.Queue
	logger  logger.Logger
}

func (r *RabbitMQConsumer) RunConsumer(ctx context.Context) {
	messages, err := r.channel.Consume(
		r.queue.Name, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		r.logger.Errorf("err consuming messages: %s", err)
	}

	for {
		select {
		case msg, ok := <-messages:
			if !ok {
				r.logger.Info("chan closed ")
			}
			r.logger.Info(string(msg.Body))

		case <-ctx.Done():
			r.logger.Info("stop consumer")
			return
		}
	}
	
}

func NewRabbitMQConsumer(config config.RabbitMQ, logger logger.Logger) (*RabbitMQConsumer, error) {

	conn, err := amqp.Dial(config.Url)
	if err != nil {
		return nil, err
	}

	r := &RabbitMQConsumer{
		conn:   conn,
		logger: logger,
	}

	err = r.Channel()
	if err != nil {
		return nil, err
	}

	err = r.Queue()
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *RabbitMQConsumer) Close() error {

	err := r.conn.Close()
	if err != nil {
		return err
	}

	err = r.channel.Close()
	if err != nil {
		return err
	}

	return nil
}

func (r *RabbitMQConsumer) Channel() error {
	if r.conn == nil {
		return fmt.Errorf("connection is nil")
	}
	ch, err := r.conn.Channel()
	if err != nil {
		return fmt.Errorf("Queue err: ", err)
	}

	r.channel = ch
	return nil

}

func (r *RabbitMQConsumer) Queue() error {

	q, err := r.channel.QueueDeclare(
		"Queue", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return fmt.Errorf("Queue err: ", err)
	}

	r.queue = &q

	return nil
}
