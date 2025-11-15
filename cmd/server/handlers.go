package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

func handlerGamelogs() func(routing.GameLog) pubsub.AckType {
	return func(gl routing.GameLog) pubsub.AckType {
		defer fmt.Printf("> ")
		err := gamelogic.WriteLog(gl)
		if err != nil {
			log.Printf("Error writing log: %v\n", err)
			log.Printf("Retrying...\n")
			return pubsub.NackRequeue
		}
		return pubsub.Ack
	}
}
