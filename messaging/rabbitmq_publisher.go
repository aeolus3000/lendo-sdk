package messaging

import (
	"bytes"
	"github.com/streadway/amqp"
)

type rabbitmqPublisher struct {
	RabbitmqAbstract
}


func NewRabbitMqPublisher(c RabbitMqConfiguration) (Publisher, error) {
	publisher := rabbitmqPublisher{}
	err := publisher.initialize(c)
	if err != nil {
		return nil, err
	}
	return &publisher, err
}

func (mq *rabbitmqPublisher) Publish(buffer bytes.Buffer) error {
	return mq.channel.Publish(
		"",           // exchange
		mq.queue.Name,       // routing key
		false,        // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  mq.configuration.ContentType,
			Body:         buffer.Bytes(),
		})
}

