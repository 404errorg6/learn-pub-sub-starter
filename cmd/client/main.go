package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	const url = "amqp://guest:guest@localhost:5672/"
	fmt.Println("Starting Peril client...")

	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalf("Error connecting to RabbitMQ: %v\n", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error creating channel: %v\n", err)
	}
	defer ch.Close()

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Error getting username: %v\n", err)
	}

	pauseQName := fmt.Sprintf("pause.%v", username)
	gameState := gamelogic.NewGameState(username)
	err = pubsub.SubscribeJSON(conn, routing.ExchangePerilDirect, pauseQName, routing.PauseKey, pubsub.Transient, handlerPause(gameState))
	if err != nil {
		log.Fatalf("Error subscribing to %v: %v\n", pauseQName, err)
	}

	moveQName := fmt.Sprintf("army_moves.%v", username)
	err = pubsub.SubscribeJSON(conn, routing.ExchangePerilTopic, moveQName, "army_moves.*", pubsub.Transient, handlerMove(gameState, ch))
	if err != nil {
		log.Fatalf("Error subscribing to %v: %v\n", moveQName, err)
	}

	warQName := fmt.Sprintf("%v.*", routing.WarRecognitionsPrefix)
	err = pubsub.SubscribeJSON(conn, routing.ExchangePerilTopic, routing.WarRecognitionsPrefix, warQName, pubsub.Durable, handlerWar(gameState, ch))
	if err != nil {
		log.Fatalf("Error subscribing to %v: %v\n", warQName, err)
	}

	startREPL(gameState, ch)
}
