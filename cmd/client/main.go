package main

import (
	"fmt"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
)
type Config struct {
	connection *amqp.Connection
	gameState *gamelogic.GameState
	Username string
}
type cliCommand struct {
	handler func(*Config,...string) error
}

func getSupportedCommands() map[string]cliCommand {
	return map[string]cliCommand {
		"spawn": cliCommand {
			handler: command_spawn,
		},
		"move": cliCommand {
			handler: command_move,
		},
		"status": cliCommand {
			handler: command_status,
		},
		"spam": cliCommand {
			handler: command_spam,
		},
		"help": cliCommand {
			handler: command_help,
		},
		"quit": cliCommand {
			handler: command_quit,
		},
	}
}

func main() {
	connectionString := "amqp://guest:guest@localhost:5672/"
	fmt.Println("Starting Peril client...")

	amqpConn, err := amqp.Dial(connectionString)
	if err != nil {
		log.Fatalf("Couldn't connect to RabbitMQ : %v", err)
	}
	defer amqpConn.Close()

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Couldn't get username during welcome: %v", err)
	}
	pauseQueue := fmt.Sprintf("pause.%s", username)
	_, _, err = pubsub.DeclareAndBind(amqpConn, routing.ExchangePerilDirect, pauseQueue, routing.PauseKey, false)
	if err != nil {
		log.Fatal("Couldn't declare and bind queue: %v", err)
	}
	config := &Config{
		connection: amqpConn,
		gameState: gamelogic.NewGameState(username),
		Username: username,
	}

	err = pubsub.SubscribeJSON(amqpConn, routing.ExchangePerilDirect, pauseQueue, routing.PauseKey, false, handlerPause(config.gameState))
	if err != nil {
		log.Fatalf("Couldn't subscribe to PauseQueue: %v", err)
	}
	
	movesQueueName := fmt.Sprintf("%s.%s", routing.ArmyMovesPrefix , username)
	movesRoutingKey := fmt.Sprintf("%s.*", routing.ArmyMovesPrefix)
	err = pubsub.SubscribeJSON(amqpConn, routing.ExchangePerilTopic, movesQueueName, movesRoutingKey, false, config.handlerMove(config.gameState))
	if err != nil {
		log.Fatalf("Couldn't subscribe to %s queue: %v", movesQueueName, err) 
	}

	warQueueName := routing.WarRecognitionsPrefix
	warRoutingKey := fmt.Sprintf("%s.*", routing.WarRecognitionsPrefix)
	err = pubsub.SubscribeJSON(amqpConn, routing.ExchangePerilTopic, warQueueName, warRoutingKey, true, config.handlerWar(config.gameState))
	if err != nil {
		log.Fatalf("Couldn't subscribe to %s queue: %v", warQueueName, err)
	}
	supportedCommands := getSupportedCommands()
	MAIN_LOOP:
	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}
		commandName := words[0]
		args := words
		command, ok := supportedCommands[commandName]
		if !ok {
			fmt.Println("Unsupported command")
			continue
		}
		err := command.handler(config, args...)
		if err != nil {
			log.Printf("Something went wrong: %v",err)
			os.Exit(1)
		}
		if commandName == "quit" {
			break MAIN_LOOP
		}

	}

}
