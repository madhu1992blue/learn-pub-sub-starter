package pubsub

import (
	"bytes"
	"encoding/json"
	"encoding/gob"
	amqp "github.com/rabbitmq/amqp091-go"
)
func SubscribeGob[T any](conn *amqp.Connection, exchange, queueName, key string, durable bool, handler func(T) Acktype) error {
	ch, _, err := DeclareAndBind(conn, exchange, queueName, key, durable)

	if err != nil {
        	return err
	}
	if err != nil {
		return err
	}
	deliveryChan, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}
	go func() {
        	defer ch.Close()
		for delivery := range deliveryChan {
			var msg T
			databuf := bytes.NewBuffer(delivery.Body)
			decoder := gob.NewDecoder(databuf)
			if err := decoder.Decode(&msg); err != nil {
				return 
			}
			acktype := handler(msg)
			if err != nil {
				return 
			}
			switch acktype {
			case Ack:
				err = delivery.Ack(false)
			case NackRequeue:
				err = delivery.Nack(false, true)
			case NackDiscard:
				err = delivery.Nack(false, false)

			}
			if err != nil {
				return
			}
		}
	}()
	return nil

}
func SubscribeJSON[T any](conn *amqp.Connection, exchange, queueName, key string, durable bool, handler func(T) Acktype) error {
	ch, _, err := DeclareAndBind(conn, exchange, queueName, key, durable)

	if err != nil {
        	return err
	}
	if err != nil {
		return err
	}
	deliveryChan, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}
	go func() {
        	defer ch.Close()
		for delivery := range deliveryChan {
			var msg T
			decoder := json.NewDecoder(bytes.NewReader(delivery.Body))
			if err := decoder.Decode(&msg); err != nil {
				return 
			}
			acktype := handler(msg)
			if err != nil {
				return 
			}
			switch acktype {
			case Ack:
				err = delivery.Ack(false)
			case NackRequeue:
				err = delivery.Nack(false, true)
			case NackDiscard:
				err = delivery.Nack(false, false)

			}
			if err != nil {
				return
			}
		}
	}()
	return nil

}
