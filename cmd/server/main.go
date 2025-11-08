package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing/pubsub"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	const connString = "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connString)
	if err != nil {
		log.Fatalf("Error occured while connecting to RabbitMQ: %v\n", err)
	}
	defer conn.Close()

	fmt.Printf("Peril game server successfully connected to RabbitMQ\n")

	amqpCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error creating channel: %v", err)
	}

	playingState := routing.PlayingState{
		IsPaused: true,
	}

	err = pubsub.PublishJson(amqpCh, routing.ExchangePerilDirect, routing.PauseKey, playingState)
	if err != nil {
		log.Fatalf("Error while publishing msg: %v", err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	defer fmt.Printf("\nRabbitMQ connection closed\n")
}
