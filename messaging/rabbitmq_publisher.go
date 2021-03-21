package messaging

import (
	"bytes"
	"errors"
	"github.com/streadway/amqp"
	"os"
	"time"
)

type rabbitmqPublisher struct {
	RabbitmqAbstract
}

func NewRabbitMqPublisher(c RabbitMqConfiguration, done chan os.Signal) (Publisher, error) {
	publisher := rabbitmqPublisher{}
	err := publisher.initialize(c, done)
	if err != nil {
		return nil, err
	}
	return &publisher, err
}

// Push will push data onto the queue, and wait for a confirmation.
// If no confirms are received until within the resendTimeout,
// it continuously resends messages until a confirmation is received.
// This will block until the server sends a confirm.
func (mq *rabbitmqPublisher) Publish(buffer bytes.Buffer) error {
	if !mq.isConnected {
		return errors.New("failed to push: not connected")
	}
	for {
		err := mq.UnsafePush(buffer)
		if err != nil {
			if err == ErrDisconnected {
				continue
			}
			return err
		}
		select {
		case confirm := <-mq.notifyConfirm:
			if confirm.Ack {
				return nil
			}
		case <-time.After(mq.configuration.ResendDelay):
		}
	}
}

// UnsafePush will push to the queue without checking for
// confirmation. It returns an error if it fails to connect.
// No guarantees are provided for whether the server will
// receive the message.
func (mq *rabbitmqPublisher) UnsafePush(buffer bytes.Buffer) error {
	if !mq.isConnected {
		return ErrDisconnected
	}
	return mq.channel.Publish(
		"",                         // Exchange
		mq.configuration.QueueName, // Routing key
		false,                      // Mandatory
		false,                      // Immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         buffer.Bytes(),
		},
	)
}
