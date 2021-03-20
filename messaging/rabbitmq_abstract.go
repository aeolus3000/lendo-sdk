package messaging

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/streadway/amqp"
)

type RabbitMqConfiguration struct {
	User string
	Password string
	Host string
	Port string
	ContentType string
	QueueName string

}

type RabbitmqAbstract struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue *amqp.Queue
	Durable bool
	AutoAck bool
	configuration RabbitMqConfiguration
}

func (mq *RabbitmqAbstract) initialize(c RabbitMqConfiguration) error {
	mq.AutoAck = false
	mq.Durable = true

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", c.User, c.Password, c.Host, c.Port))
	if err != nil {
		return err
	}
	mq.conn = conn

	channel, err := conn.Channel()
	if err != nil {
		err = multierror.Append(err, mq.shutdownConnection())
		return err
	}
	mq.channel = channel

	queue, err := channel.QueueDeclare(
		c.QueueName, // name
		mq.Durable,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		err = multierror.Append(err, mq.shutdown())
		return err
	}
	mq.queue = &queue

	mq.configuration = c

	return nil
}

func (mq *RabbitmqAbstract) shutdown() error {
	var err error = nil
	connectionCloseErr := mq.shutdownConnection()
	if connectionCloseErr != nil {
		err = multierror.Append(err, connectionCloseErr)
	}

	channelCloseErr := mq.shutdownChannel()
	if channelCloseErr != nil && err == nil {
		err = multierror.Append(err, connectionCloseErr)
	}

	return err
}

func (mq *RabbitmqAbstract) shutdownConnection() error {
	connectionCloseErr := mq.conn.Close()
	return connectionCloseErr
}

func (mq *RabbitmqAbstract) shutdownChannel() error {
	channelCloseErr := mq.channel.Close()
	return channelCloseErr
}
