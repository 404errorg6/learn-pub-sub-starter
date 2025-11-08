package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	const url = "amqp://guest:guest@localhost:5672/"
	fmt.Println("Starting Peril client...")
	amqpConn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalf("Error connecting to RabbitMQ: %v\n", err)
	}
	defer amqpConn.Close()

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Error getting username: %v\n", err)
	}

	queueName := fmt.Sprintf("pause.%v", username)
	_, _, err = pubsub.DeclareAndBind(amqpConn, routing.ExchangePerilDirect, routing.PauseKey, queueName, pubsub.Transient)
	if err != nil {
		log.Fatalf("Error binding queue: %v\n", err)
	}

	// make interactive
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	<-signalCh
	fmt.Printf("Exiting game...\n")
}
