package messaging

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
	"time"
)

const (
	// When reconnecting to the server after connection failure
	reconnectDelay = 5 * time.Second
)

var (
	ErrDisconnected = errors.New("disconnected from rabbitmq, trying to reconnect")
)

type RabbitMqConfiguration struct {
	User        string        `default:"guest"`
	Password    string        `default:"guest"`
	Host        string        `required:"true"`
	Port        string        `default:"5672"`
	ContentType string        `default:"foobar"`
	QueueName   string        `default:"foobar"`
	ResendDelay time.Duration `default:"1s"`
}

type RabbitmqAbstract struct {
	connection    *amqp.Connection
	channel       *amqp.Channel
	done          <-chan os.Signal
	notifyClose   chan *amqp.Error
	notifyConfirm chan amqp.Confirmation
	isConnected   bool
	alive         bool
	configuration RabbitMqConfiguration
}

func (mq *RabbitmqAbstract) initialize(c RabbitMqConfiguration, done <-chan os.Signal) error {

	mq.configuration = c
	mq.alive = true
	mq.done = done

	go mq.handleReconnect(c)

	return nil
}

// handleReconnect will wait for a connection error on
// notifyClose, and then continuously attempt to reconnect.
func (mq *RabbitmqAbstract) handleReconnect(c RabbitMqConfiguration) {
	for mq.alive {
		mq.isConnected = false
		t := time.Now()
		fmt.Printf("Attempting to connect to rabbitMQ: %s\n", createAddress(c))
		var retryCount int
		for !mq.connect() {
			if !mq.alive {
				return
			}
			select {
			case <-mq.done:
				return
			case <-time.After(reconnectDelay + time.Duration(retryCount)*time.Second):
				log.Printf("disconnected from rabbitMQ and failed to connect")
				retryCount++
			}
		}
		log.Printf("Connected to rabbitMQ in: %vms", time.Since(t).Milliseconds())
		select {
		case <-mq.done:
			return
		case <-mq.notifyClose:
		}
	}
}

// connect will make a single attempt to connect to
// RabbitMq. It returns the success of the attempt.
func (mq *RabbitmqAbstract) connect() bool {
	conn, err := amqp.Dial(createAddress(mq.configuration))
	if err != nil {
		return false
	}
	ch, err := conn.Channel()
	if err != nil {
		return false
	}
	ch.Confirm(false)
	_, err = ch.QueueDeclare(
		mq.configuration.QueueName,
		true,  // Durable
		false, // Delete when unused
		false, // Exclusive
		false, // No-wait
		nil,   // Arguments
	)
	if err != nil {
		return false
	}

	mq.changeConnection(conn, ch)
	mq.isConnected = true
	return true
}

// changeConnection takes a new connection to the queue,
// and updates the channel listeners to reflect this.
func (c *RabbitmqAbstract) changeConnection(connection *amqp.Connection, channel *amqp.Channel) {
	c.connection = connection
	c.channel = channel
	c.notifyClose = make(chan *amqp.Error)
	c.notifyConfirm = make(chan amqp.Confirmation)
	c.channel.NotifyClose(c.notifyClose)
	c.channel.NotifyPublish(c.notifyConfirm)
}

// Close will cleanly shutdown the channel and connection after there are no messages in the system.
func (c *RabbitmqAbstract) Close() error {
	if !c.isConnected {
		return nil
	}
	c.alive = false
	fmt.Println("Waiting for current messages to be processed...")

	fmt.Println("Closing consumer: ", 1)
	errCancel := c.channel.Cancel(consumerName(1), false)
	if errCancel != nil {
		return fmt.Errorf("error canceling consumer %s: %v", consumerName(1), errCancel)
	}

	err := c.channel.Close()
	if err != nil {
		return err
	}
	err = c.connection.Close()
	if err != nil {
		return err
	}
	c.isConnected = false
	fmt.Println("gracefully stopped rabbitMQ connection")
	return nil
}

func createAddress(c RabbitMqConfiguration) string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s", c.User, c.Password, c.Host, c.Port)
}

func consumerName(i int) string {
	return fmt.Sprintf("go-consumer-%v", i)
}
