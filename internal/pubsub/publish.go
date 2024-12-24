package pubsub

import (
	"context"
	"bytes"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"encoding/gob"
)


func PublishGob[T any](ch *amqp.Channel, exchange, key string, val T) error {

	var databuf bytes.Buffer
	encoder := gob.NewEncoder(&databuf)
	if err := encoder.Encode(val); err != nil {
		return err
	}
	
	err := ch.PublishWithContext(
		context.Background(),
		exchange,
		key,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/gob",
			Body:        databuf.Bytes(),
		},
	)
	if err != nil {
		return err
	}
	return nil
}
func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {

	bodyBytes, err := json.Marshal(val)
	if err != nil {
		return err
	}
	err = ch.PublishWithContext(
		context.Background(),
		exchange,
		key,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bodyBytes,
		},
	)
	if err != nil {
		return err
	}
	return nil
}
