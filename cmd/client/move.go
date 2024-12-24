package main

import "github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
import "github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
import "fmt"

func command_move(config *Config, args ...string) error {
	armyMove, err := config.gameState.CommandMove(args)
	movePubAMQPChan, err := config.connection.Channel()
	if err != nil {
		return err
	}
	moveRoutingKey := fmt.Sprintf("%s.%s", routing.ArmyMovesPrefix, config.Username)
	pubsub.PublishJSON(movePubAMQPChan, routing.ExchangePerilTopic, moveRoutingKey, armyMove)
	if err != nil {
		return err
	}
	fmt.Println("The move message was published successfully")
	return nil
}
