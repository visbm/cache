package rabbitmq

import (
	"cache/internal/config"
	"fmt"
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   *amqp.Queue
}

func NewRabbitMQ(config config.RabbitMQ) (*RabbitMQ, error) {

	conn, err := amqp.Dial(config.Url)
	if err != nil {
		return nil, err
	}

	r := &RabbitMQ{
		conn: conn,
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

func (r *RabbitMQ) Close() error {

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

func (r *RabbitMQ) Channel() ( error) {
	if r.conn == nil {
		return  errors.New("Connection is nil")
	}
	ch, err := r.conn.Channel()
	if err != nil {
		return fmt.Errorf("Queue err:%s ", err)
	}

	r.channel = ch 
	return nil

}

func (r *RabbitMQ) Queue() error {
	
	q, err := r.channel.QueueDeclare(
		"Queue", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return fmt.Errorf("Queue err:%s ", err)
	}

	r.queue = &q

	return nil
}
