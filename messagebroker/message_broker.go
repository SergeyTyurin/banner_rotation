package messagebroker

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/SergeyTyurin/banner-rotation/configs"
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
	url := config.URL()
	url = strings.ReplaceAll(url, "{host}", config.Host())
	url = strings.ReplaceAll(url, "{port}", strconv.Itoa(config.Port()))
	url = strings.ReplaceAll(url, "{user}", os.Getenv("MQ_USER"))
	url = strings.ReplaceAll(url, "{password}", os.Getenv("MQ_PASSWORD"))

	conn, err := amqp.Dial(url)
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
