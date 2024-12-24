package main
import (
	"time"
	"fmt"
	"strconv"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"	
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"	
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"

)
func command_spam(c *Config, words ...string) error {
	if len(words) < 2 {
		fmt.Println("Need an argument for how many spam messages")
		fmt.Print("> ")
		return nil
	}
	numSpam , err := strconv.Atoi(words[1])
	if err != nil {
		return err
	}
	amqpChan, err := c.connection.Channel()
	if err != nil {
		return err
	}
	routingKey := routing.GameLogSlug+"."+c.Username
	for i:=0; i <numSpam ; i++ {
		err = pubsub.PublishGob(
        		amqpChan, routing.ExchangePerilTopic, routingKey,
			routing.GameLog {
        			CurrentTime: time.Now(),
				Username: c.Username,
				Message: gamelogic.GetMaliciousLog(),
			},
		)
		if err != nil {
			return err
		}
	}
	return nil
}
