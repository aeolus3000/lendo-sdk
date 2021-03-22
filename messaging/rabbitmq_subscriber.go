package messaging

import (
	"bytes"
	"github.com/hashicorp/go-multierror"
	"github.com/streadway/amqp"
	"os"
	"time"
)

type rabbitmqSubscriber struct {
	RabbitmqAbstract
}

func NewRabbitMqSubscriber(c RabbitMqConfiguration, done <-chan os.Signal) (Subscriber, error) {
	subscriber := rabbitmqSubscriber{}
	errInit := subscriber.initialize(c, done)
	if errInit != nil {
		return nil, errInit
	}
	return &subscriber, nil
}

func (mq *rabbitmqSubscriber) Consume() (<-chan *Message, error) {
	for {
		if mq.isConnected {
			break
		}
		time.Sleep(1 * time.Second)
	}
	errQos := mq.channel.Qos(
		mq.configuration.PrefetchCount,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if errQos != nil {
		return nil, multierror.Append(errQos, mq.Close())
	}
	sourceMsgs, err := mq.channel.Consume(
		mq.configuration.QueueName, // queue
		consumerName(1),            // consumer
		false,                      // auto-ack
		false,                      // exclusive
		false,                      // no-local
		false,                      // no-wait
		nil,                        // args
	)
	if err != nil {
		return nil, err
	}
	destinationMsgs := make(chan *Message, cap(sourceMsgs))
	go forwardMessages(sourceMsgs, destinationMsgs)
	return destinationMsgs, nil
}

func forwardMessages(source <-chan amqp.Delivery, destination chan *Message) {
	defer close(destination)
	for d := range source {
		buffer := bytes.Buffer{}

		n, err := buffer.Write(d.Body)
		if err != nil || n != len(d.Body) {
			//log error here
			continue
		}
		message := Message{
			Body:           buffer,
			Acknowledge:    acknowledgeMessage(d),
			NotAcknowledge: notAcknowledgeMessage(d),
			Reject:         rejectMessage(d),
		}
		destination <- &message
	}
}

func acknowledgeMessage(source amqp.Delivery) func() error {
	return func() error {
		return source.Ack(false)
	}
}

func notAcknowledgeMessage(source amqp.Delivery) func(bool) error {
	return func(requeue bool) error {
		return source.Nack(false, requeue)
	}
}

func rejectMessage(source amqp.Delivery) func(bool) error {
	return func(requeue bool) error {
		return source.Reject(requeue)
	}
}
