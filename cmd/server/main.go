package main

import (
	"fmt"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"os/signal"
)

type Config struct {
	amqpConn *amqp.Connection
}
type cliCommand struct {
	handler func(c *Config, args ...string)
}

func getSupportedCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"pause": cliCommand{
			handler: command_pause,
		},
		"resume": cliCommand{
			handler: command_resume,
		},
	}
}
func main() {
	connectionString := "amqp://guest:guest@localhost:5672/"
	fmt.Println("Starting Peril server...")
	amqpConn, err := amqp.Dial(connectionString)
	if err != nil {
		log.Fatalf("Couldn't connect to RabbitMQ server: %v\n", err)
	}
	defer amqpConn.Close()

	config := &Config{
		amqpConn: amqpConn,
	}
	supportedCommands := getSupportedCommands()

	fmt.Println("AMQP connection was successful")
	
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
			
	gamelogic.PrintServerHelp()
	MAIN_LOOP:
	for {
		select {
		case <-sigChan:
			<-sigChan
			fmt.Println("Shutting down Peril server...")
			break MAIN_LOOP
		default:
			words := gamelogic.GetInput()
			if len(words) == 0 {
				continue
			}
			commandName := words[0]
			if commandName == "quit" {
				fmt.Println("Exiting the program")
				break MAIN_LOOP
			}
			command, ok := supportedCommands[commandName]
			if !ok {
				fmt.Println("Unknown command")
				continue
			}
			command.handler(config, words[1:]...)
		}
	}
}
