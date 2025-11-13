package main

import (
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"github.com/rabbitmq/amqp091-go"
)

func handleCmds(ch *amqp091.Channel) {
	var playingState routing.PlayingState
	for {
		input := gamelogic.GetInput()
		if input == nil {
			continue
		}

		cmd := input[0]

		if cmd == "pause" {
			log.Printf("Pausing the game...\n")
			playingState.IsPaused = true
			err := pubsub.PublishJson(ch, routing.ExchangePerilDirect, routing.PauseKey, playingState)
			if err != nil {
				log.Fatalf("Error publishing pause msg: %v", err)
			}
			continue
		}

		if cmd == "resume" {
			log.Printf("Resuming the game...\n")
			playingState.IsPaused = false
			err := pubsub.PublishJson(ch, routing.ExchangePerilDirect, routing.PauseKey, playingState)
			if err != nil {
				log.Fatalf("Error publishing resumme msg: %v", err)
			}
			continue
		}

		if cmd == "quit" {
			log.Printf("Exiting the game\n")
			return
		}

		log.Printf("Unknown command: %v\n", cmd)
	}
}
