package main

import (
	"time"
	"log"
	"fmt"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) pubsub.Acktype {
	return func(ps routing.PlayingState) pubsub.Acktype {
		defer fmt.Print("> ")
		gs.HandlePause(ps)
		return pubsub.Ack
	}
}

func (c *Config) handlerMove(gs *gamelogic.GameState) func(gamelogic.ArmyMove) pubsub.Acktype {
	return func(am gamelogic.ArmyMove) pubsub.Acktype {
		defer fmt.Print("> ")
		moveOutcome := gs.HandleMove(am)
		switch moveOutcome {
		case gamelogic.MoveOutComeSafe:
			return pubsub.Ack
		case gamelogic.MoveOutcomeMakeWar:
			routingKey := routing.WarRecognitionsPrefix+"."+c.Username
			amqpChan, err := c.connection.Channel()
			if err != nil {
				log.Println("Error creating channel for publishing RecognitionOfWar")
				return pubsub.NackRequeue
			}
			err = pubsub.PublishJSON(
				amqpChan, routing.ExchangePerilTopic, routingKey, 
				gamelogic.RecognitionOfWar{
					Attacker: am.Player,
					Defender: gs.Player,
				},
			)
			if err != nil {
				log.Printf("Error publishing: %v\n", err)
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		case gamelogic.MoveOutcomeSamePlayer:
			return pubsub.NackDiscard
		default:
			return pubsub.NackDiscard
		}
	}
}


func (c *Config) handlerLog(_ *gamelogic.GameState) func(routing.GameLog) pubsub.Acktype {
	return func(gl routing.GameLog) pubsub.Acktype{
		defer fmt.Print("> ")
		err := gamelogic.WriteLog(gl)
		if err != nil {
			return pubsub.NackRequeue
		}
		return pubsub.Ack
	}

}

func (c *Config) handlerWar(gs *gamelogic.GameState) func(gamelogic.RecognitionOfWar) pubsub.Acktype {

	return func(rw gamelogic.RecognitionOfWar) pubsub.Acktype {
		defer fmt.Print("> ")
		outcome, winner, loser := gs.HandleWar(rw)

		routingKey := routing.GameLogSlug+"."+c.Username
		amqpChan, err := c.connection.Channel()
		if err != nil {
			log.Println("Error creating channel for publishing RecognitionOfWar")
				return pubsub.NackRequeue
		}
		
		switch outcome {
		case gamelogic.WarOutcomeNotInvolved:
			return pubsub.NackRequeue
		case gamelogic.WarOutcomeNoUnits:
			return pubsub.NackDiscard
		case gamelogic.WarOutcomeOpponentWon:
			fallthrough
		case gamelogic.WarOutcomeYouWon:
			err = pubsub.PublishGob(
				amqpChan, routing.ExchangePerilTopic, routingKey, 
				routing.GameLog {
					CurrentTime: time.Now(),
					Username: c.Username,
					Message: fmt.Sprintf("%s won a war against %s", winner, loser), 
				},
			)
			if err != nil {
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		case gamelogic.WarOutcomeDraw:
			err = pubsub.PublishGob(
				amqpChan, routing.ExchangePerilTopic, routingKey, 
				routing.GameLog {
					CurrentTime: time.Now(),
					Username: c.Username,
					Message: fmt.Sprintf("A war between %s and %s resulted in a draw", winner, loser),
				},
			)
			if err != nil {
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		default:
			log.Println("Unexpected outcome")
			return pubsub.NackDiscard
		}

	}
}
