package pubsub

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType int

const (
	Durable SimpleQueueType = iota
	Transient
)

type AckType int

const (
	Ack AckType = iota
	NackRequeue
	NackDiscard
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
	handler func(T) AckType,
) error {
	ch, _, err := DeclareAndBind(conn, exchange, key, queueName, queueType)
	if err != nil {
		return err
	}

	deliveryCh, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	var value T

	go func(value *T) {
		for msg := range deliveryCh {
			err := json.Unmarshal(msg.Body, value)
			if err != nil {
				log.Printf("Error unmarshalling data: %v\n", err)
			}

			ack := handler(*value)
			switch ack {
			case Ack:
				err = msg.Ack(false)
				if err != nil {
					log.Printf("Error acknowledging msg: %v\n", err)
				}
				fmt.Printf("(Ack used)\n")

			case NackRequeue:
				err = msg.Nack(false, true)
				if err != nil {
					log.Printf("Error NackRequeue-ing msg: %v\n", err)
				}
				fmt.Printf("(NackRequeue used)\n")

			case NackDiscard:
				err = msg.Nack(false, false)
				if err != nil {
					log.Printf("Error NackDiscard-ing msg: %v\n", err)
				}
				fmt.Printf("(NackDiscard used)\n")
			}
			fmt.Printf("> ")
		}
	}(&value)

	return err
}

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	key,
	queueName string,
	queueType SimpleQueueType,
) (*amqp.Channel, amqp.Queue, error) {
	var queue amqp.Queue
	ch, err := conn.Channel()
	if err != nil {
		return ch, queue, err
	}

	var durable, autoDel, exclusive bool
	switch queueType {
	case Durable:
		durable = true

	case Transient:
		autoDel = true
		exclusive = true
	}

	args := make(amqp.Table)
	args["x-dead-letter-exchange"] = "peril_dlx"
	queue, err = ch.QueueDeclare(queueName, durable, autoDel, exclusive, false, args)
	if err != nil {
		return ch, queue, err
	}

	err = ch.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		return ch, queue, err
	}
	return ch, queue, nil
}
