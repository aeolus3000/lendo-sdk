// +build integration_test

//This test need the following docker container running
//run -d --hostname my-rabbit6 --name some-rabbit6 --network host rabbitmq:3.8.14-management

package messaging

import (
	"bytes"
	"fmt"
	"github.com/aeolus3000/lendo-sdk/banking"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"log"
	"os"
	"os/signal"
	"sync/atomic"

	"syscall"
	"testing"
	"time"
)

var (
	testPubSub = TestPubSub{}
	config     = RabbitMqConfiguration{
		User:        "guest",
		Password:    "guest",
		Host:        "localhost",
		Port:        "5672",
		ContentType: "text/plain",
		QueueName:   "myqueue",
		ResendDelay: 1 * time.Second,
	}
)

type TestPubSub struct {
	pub Publisher
	sub Subscriber
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	publisher, _ := NewRabbitMqPublisher(config, sigs)
	subscriber, _ := NewRabbitMqSubscriber(config, sigs)

	testPubSub = TestPubSub{
		pub: publisher,
		sub: subscriber,
	}

	rabbitMqConnectTime := 2 * time.Second
	time.Sleep(rabbitMqConnectTime)
}

func shutdown() {
	testPubSub.pub.Close()
	testPubSub.sub.Close()
}

func TestPolling(t *testing.T) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	config2     := RabbitMqConfiguration{
		User:        "guest",
		Password:    "guest",
		Host:        "localhost",
		Port:        "5672",
		ContentType: "text/plain",
		QueueName:   "topolling",
		ResendDelay: 1 * time.Second,
	}
	publisher, _ := NewRabbitMqPublisher(config2, sigs)
	time.Sleep(2 * time.Second)
	application := banking.Application{
		Id:        uuid.NewString(),
	}
	bytesApplication, _ := proto.Marshal(&application)
	bytesApplication2 := bytes.NewBuffer(bytesApplication)
	err := publisher.Publish(bytesApplication2)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestRabbitMq(t *testing.T) {
	//Given expected string
	expectedString := "Hallo Welt"
	//When publishing expected string
	buffer := bytes.Buffer{}
	buffer.Write([]byte(expectedString))
	err := testPubSub.pub.Publish(&buffer)
	if err != nil {
		t.Errorf("Could not publish to rabbitmq; error: %v", err)
	}
	//And subscribing on the same channel
	msgs, _ := testPubSub.sub.Consume()
	//And an atomic counter
	var ops uint64 = 1
	//then
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body.String())
			log.Printf("Done")
			if d.Body.String() != expectedString {
				t.Errorf("Got: %s; Want: %s", d.Body.String(), expectedString)
			}
			if ops != 1 {
				t.Errorf("The subscriber received unexpectedly more than one message")
			}
			_ = d.Acknowledge()
			atomic.AddUint64(&ops, 1)
		}
	}()

	select {
	case <-time.After(5 * time.Second):
	}

}
