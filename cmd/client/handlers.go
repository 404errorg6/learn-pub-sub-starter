package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"github.com/rabbitmq/amqp091-go"
)

func handlerWar(gs *gamelogic.GameState, ch *amqp091.Channel) func(gamelogic.RecognitionOfWar) pubsub.AckType {
	defer fmt.Printf("> ")
	return func(row gamelogic.RecognitionOfWar) pubsub.AckType {

		var msg string
		key := fmt.Sprintf("%v.%v", routing.GameLogSlug, row.Attacker.Username)
		gl := routing.GameLog{}
		outcome, winner, loser := gs.HandleWar(row)

		switch outcome {
		case gamelogic.WarOutcomeNotInvolved:
			return pubsub.NackRequeue

		case gamelogic.WarOutcomeNoUnits:
			return pubsub.NackDiscard

		case gamelogic.WarOutcomeOpponentWon:
		case gamelogic.WarOutcomeYouWon:
			msg = fmt.Sprintf("%v won a war against %v\n", winner, loser)
			gl.Message = msg

			err := pubsub.PublishGob(ch, routing.ExchangePerilTopic, key, gl)
			if err != nil {
				return pubsub.NackRequeue
			}
			return pubsub.Ack

		case gamelogic.WarOutcomeDraw:
			msg = fmt.Sprintf("A between %v and %v resulted in a draw\n", winner, loser)
			gl.Message = msg

			err := pubsub.PublishGob(ch, routing.ExchangePerilTopic, key, gl)
			if err != nil {
				return pubsub.NackRequeue
			}
			return pubsub.Ack

		default:
			msg = fmt.Sprintf("Invalid outcome: %v\n", outcome)
		}
		return pubsub.Ack
	}
}

func handlerMove(gs *gamelogic.GameState, ch *amqp091.Channel) func(gamelogic.ArmyMove) pubsub.AckType {
	return func(move gamelogic.ArmyMove) pubsub.AckType {
		outcome := gs.HandleMove(move)
		defer fmt.Printf("> ")
		switch outcome {
		case gamelogic.MoveOutComeSafe:
			return pubsub.Ack

		case gamelogic.MoveOutcomeMakeWar:
			key := fmt.Sprintf("%v.%v", routing.WarRecognitionsPrefix, gs.Player.Username)
			rw := gamelogic.RecognitionOfWar{
				Attacker: move.Player,
				Defender: gs.Player,
			}

			err := pubsub.PublishJson(ch, routing.ExchangePerilTopic, key, rw)
			if err != nil {
				fmt.Printf("Error publishing move: %v\n", err)
				return pubsub.NackRequeue
			}
			return pubsub.Ack

		}
		return pubsub.NackDiscard
	}
}

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) pubsub.AckType {
	return func(ps routing.PlayingState) pubsub.AckType {
		gs.HandlePause(ps)
		defer fmt.Printf("> ")
		return pubsub.Ack
	}

}
