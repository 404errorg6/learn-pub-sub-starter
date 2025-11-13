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
	fmt.Println("Starting Peril server...")
	const url = "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalf("Error occured while connecting to RabbitMQ: %v\n", err)
	}
	defer conn.Close()

	fmt.Printf("Peril game server successfully connected to RabbitMQ\n")

	amqpCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error creating channel: %v", err)
	}

	topicKey := "game_logs.*"
	_, _, err = pubsub.DeclareAndBind(conn, routing.ExchangePerilTopic, topicKey, routing.GameLogSlug, pubsub.Durable)
	if err != nil {
		log.Fatalf("Error while binding server to queue: %v\n", err)
	}

	gamelogic.PrintServerHelp()

	handleCmds(amqpCh)
}
