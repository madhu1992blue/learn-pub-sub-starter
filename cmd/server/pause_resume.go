package main

import (
	"fmt"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"log"
)

func command_pause(config *Config, args ...string) {
	mqChan, err := config.amqpConn.Channel()
	if err != nil {
		log.Fatalf("Couldn't create AMQP channel: %v", err)
	}
	fmt.Println("Publishing a 'pause' message to the queue")
	pubsub.PublishJSON(mqChan, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
		IsPaused: true,
	})
}

func command_resume(config *Config, args ...string) {
	mqChan, err := config.amqpConn.Channel()
	if err != nil {
		log.Fatalf("Couldn't create AMQP channel: %v", err)
	}
	fmt.Println("Publishing a 'resume' message to the queue")
	pubsub.PublishJSON(mqChan, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
		IsPaused: false,
	})
}
