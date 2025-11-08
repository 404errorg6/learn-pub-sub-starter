package pubsub

import amqp "github.com/rabbitmq/amqp091-go"

type SimpleQueueType int

const (
	Durable SimpleQueueType = iota
	Transient
)

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
	if queueType == Durable {
		durable = true
	} else if queueType == Transient {
		autoDel = true
		exclusive = true
	}

	queue, err = ch.QueueDeclare(queueName, durable, autoDel, exclusive, false, nil)
	if err != nil {
		return ch, queue, err
	}

	err = ch.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		return ch, queue, err
	}
	return ch, queue, nil
}
