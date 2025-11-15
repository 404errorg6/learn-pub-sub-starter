package pubsub

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"

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

func SubscribeGob[T any](conn *amqp.Connection, exchange, queueName, key string, queueType SimpleQueueType, handler func(T) AckType) error {
	err := subscribe(
		conn,
		exchange,
		queueName,
		key,
		queueType,
		handler,
		func(body []byte) (T, error) {
			var converted T
			buffer := bytes.NewReader(body)
			decoder := gob.NewDecoder(buffer)

			err := decoder.Decode(&converted)
			if err != nil {
				return converted, err
			}
			return converted, nil
		})
	return err
}

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
	handler func(T) AckType,
) error {
	err := subscribe(
		conn,
		exchange,
		queueName,
		key,
		queueType,
		handler,
		func(body []byte) (T, error) {
			var val T
			err := json.Unmarshal(body, &val)

			if err != nil {
				return val, err
			}
			return val, nil
		})
	return err
}

func subscribe[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
	handler func(T) AckType,
	unmarshaller func([]byte) (T, error),
) error {
	ch, _, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return err
	}

	err = ch.Qos(10, 0, false)
	if err != nil {
		return err
	}

	deliveryCh, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for msg := range deliveryCh {
			val, err := unmarshaller(msg.Body)
			if err != nil {
				fmt.Printf("Error occured while unmarshalling msg: %v\n", err)
				continue
			}

			ackType := handler(val)
			switch ackType {
			case Ack:
				err = msg.Ack(false)
				if err != nil {
					fmt.Printf("Error acknowledging msg: %v\n", err)
					continue
				}
				fmt.Printf("(Ack used)\n")

			case NackRequeue:
				err = msg.Nack(false, true)
				if err != nil {
					fmt.Printf("Error NackRequeue-ing msg: %v\n", err)
				}
				fmt.Printf("(NackRequeue used)\n")

			case NackDiscard:
				err = msg.Nack(false, false)
				if err != nil {
					fmt.Printf("Error NackDiscard-ing msg: %v\n", err)
				}
				fmt.Printf("(NackDiscard used)\n")
			}
			fmt.Printf("> ")
		}
	}()
	return nil
}

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
) (*amqp.Channel, amqp.Queue, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("could not create channel: %v", err)
	}

	queue, err := ch.QueueDeclare(
		queueName,            // name
		queueType == Durable, // durable
		queueType != Durable, // delete when unused
		queueType != Durable, // exclusive
		false,                // no-wait
		amqp.Table{
			"x-dead-letter-exchange": "peril_dlx",
		},
	)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("could not declare queue: %v", err)
	}

	err = ch.QueueBind(
		queue.Name, // queue name
		key,        // routing key
		exchange,   // exchange
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("could not bind queue: %v", err)
	}
	return ch, queue, nil
}
