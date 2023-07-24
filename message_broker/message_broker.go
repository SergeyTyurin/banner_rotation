package message_broker

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/SergeyTyurin/banner_rotation/configs"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageBroker interface {
	Connect(configs.MessageBrokerConfig) (func(), error)
	SendRegisterTransitionEvent(string) error
	SendSelectFromRotationEvent(string) error

	GetRegisterTransitionEvent() (string, error)
	GetSelectFromRotationEvent() (string, error)
}

type messageBrokerImpl struct {
	conn          *amqp.Connection
	ch            *amqp.Channel
	registerQueue amqp.Queue
	selectQueue   amqp.Queue
}

func NewBroker() MessageBroker {
	return &messageBrokerImpl{}
}

func (m *messageBrokerImpl) Connect(config configs.MessageBrokerConfig) (func(), error) {
	connStr := fmt.Sprintf("amqp://%s:%s@%s:%d/",
		os.Getenv("MQ_USER"),
		os.Getenv("MQ_PASSWORD"),
		config.Host(),
		config.Port())

	conn, err := amqp.Dial(connStr)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	rq, err := ch.QueueDeclare(
		"RegisterTransition", // name
		false,                // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	sq, err := ch.QueueDeclare(
		"SelectFromRotation", // name
		false,                // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	m.conn = conn
	m.ch = ch
	m.registerQueue = rq
	m.selectQueue = sq

	return func() {
		m.ch.Close()
		m.conn.Close()
	}, nil
}

func (m *messageBrokerImpl) SendRegisterTransitionEvent(body string) error {
	err := m.ch.PublishWithContext(context.Background(),
		"",                   // exchange
		m.registerQueue.Name, // routing key
		false,                // mandatory
		false,                // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		return err
	}
	log.Printf(" [x] Sent %s\n", body)
	return nil
}

func (m *messageBrokerImpl) SendSelectFromRotationEvent(body string) error {
	err := m.ch.PublishWithContext(context.Background(),
		"",                 // exchange
		m.selectQueue.Name, // routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		return err
	}
	log.Printf(" [x] Sent %s\n", body)
	return nil
}

func (m *messageBrokerImpl) GetSelectFromRotationEvent() (string, error) {
	msgs, err := m.ch.Consume(
		m.selectQueue.Name, // queue
		"",                 // consumer
		true,               // auto-ack
		false,              // exclusive
		false,              // no-local
		false,              // no-wait
		nil,                // args
	)
	if err != nil {
		return "", err
	}
	for d := range msgs {
		return string(d.Body), nil
	}
	return "", nil
}

func (m *messageBrokerImpl) GetRegisterTransitionEvent() (string, error) {
	msgs, err := m.ch.Consume(
		m.registerQueue.Name, // queue
		"",                   // consumer
		true,                 // auto-ack
		false,                // exclusive
		false,                // no-local
		false,                // no-wait
		nil,                  // args
	)
	if err != nil {
		return "", err
	}
	for d := range msgs {
		return string(d.Body), nil
	}
	return "", nil
}
