package messaging

import (
	"bytes"
	"log"
	"testing"
	"time"
)

func TestRabbitMq(t *testing.T) {
	config := RabbitMqConfiguration{
		User:        "guest",
		Password:    "guest",
		Host:        "localhost",
		Port:        "5672",
		ContentType: "text/plain",
		QueueName:   "myqueue",
	}
	
	publisher, _ := NewRabbitMqPublisher(config)
	subscriber, _ := NewRabbitMqSubscriber(config)

	msgs, _ := subscriber.Consume()

	buffer := bytes.Buffer{}
	buffer.Write([]byte("Hallo Welt"))
	_ = publisher.Publish(buffer)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body.String())
			log.Printf("Done")
			_ = d.Acknowledge()
		}
	}()

	select {
	case <-time.After(5 * time.Second):
	}
}