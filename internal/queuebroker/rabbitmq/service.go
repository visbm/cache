package rabbitmq

import (
	"cache/pkg/logger"
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQService struct {
	cl     *RabbitMQ
	logger logger.Logger
}

func NewRabbitMQService(cl *RabbitMQ, logger logger.Logger) *RabbitMQService {
	return &RabbitMQService{cl: cl, logger: logger}

}

func (r *RabbitMQService) SendMessage(message string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.cl.channel.PublishWithContext(ctx,
		"",              // exchange
		r.cl.queue.Name, // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})

	if err != nil {
		return err
	}

	err = r.cl.channel.Confirm(false)
	if err != nil {
		return err
	}	
	return nil

}
