package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"github.com/rabbitmq/amqp091-go"
)

func startREPL(gameState *gamelogic.GameState, ch *amqp091.Channel) {
	for {

		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}
		cmd := input[0]

		if cmd == "spawn" {
			err := gameState.CommandSpawn(input)
			if err != nil {
				fmt.Printf("Error occured while spawning: %v\n", err)
			}
			continue
		}

		if cmd == "move" {
			move, err := gameState.CommandMove(input)
			if err != nil {
				fmt.Printf("Error occured while moving army: %v\n", err)
				continue
			}

			key := fmt.Sprintf("army_moves.%v", gameState.Player.Username)
			err = pubsub.PublishJson(ch, routing.ExchangePerilTopic, key, move)
			if err != nil {
				fmt.Printf("Error occured while publishing move: %v\n", err)
			}

			continue
		}

		if cmd == "status" {
			gameState.CommandStatus()
			continue
		}

		if cmd == "help" {
			gamelogic.PrintClientHelp()
			continue
		}

		if cmd == "spam" {
			fmt.Printf("Spamming not allowed yet!\n")
			continue
		}

		if cmd == "quit" {
			gamelogic.PrintQuit()
			break
		}

		fmt.Printf("\"%v\" command not found\n", cmd)
	}
}
