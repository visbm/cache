package queuebroker

import (
	"cache/internal/queuebroker/rabbitmq"
	"cache/pkg/logger"
)

type Broker interface {
	SendMessage(message string) error
}

func NewBrokerService(r *rabbitmq.RabbitMQ, logger logger.Logger) Broker {
	return rabbitmq.NewRabbitMQService(r, logger)
}
