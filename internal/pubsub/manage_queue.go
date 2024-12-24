package pubsub

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func DeclareAndBind(conn *amqp.Connection, exchange, queueName, routingKey string, durable bool) (*amqp.Channel, amqp.Queue, error) {

	amqpChan, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	exclusive, autoDelete := true, true
	if durable {
		exclusive, autoDelete = false, false
	}

	queue, err := amqpChan.QueueDeclare(queueName, durable, autoDelete, exclusive, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, err
	}
	err = amqpChan.QueueBind(queueName, routingKey, exchange, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, err
	}
	return amqpChan, queue, nil
}
