package rabbitmq

import (
	"context"
	"fmt"
	"log"

	"github.com/aclgo/simple-api-gateway/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	QueueName = "simple-api-gateway"
)

type Rabbitmq struct {
	conn           *amqp.Connection
	ch             *amqp.Channel
	queueProducer  *amqp.Queue
	notifyReturnCh chan amqp.Return
	confirmationCh chan amqp.Confirmation
}

func NewRabbitMq(cfg *config.Config) *Rabbitmq {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatalf("amqp.Dial: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("faile to open channel: %v", err)
	}

	if err := ch.Confirm(false); err != nil {
		log.Fatalf("failed to put channel confirm\n")
	}

	q, err := ch.QueueDeclare(
		QueueName,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatalf("Failed to declare queue: %v\n", err)
	}

	return &Rabbitmq{
		conn:           conn,
		ch:             ch,
		queueProducer:  &q,
		notifyReturnCh: ch.NotifyReturn(make(chan amqp.Return)),
		confirmationCh: ch.NotifyPublish(make(chan amqp.Confirmation)),
	}
}

func (r *Rabbitmq) Producer(ctx context.Context, message string) error {

	err := r.ch.PublishWithContext(ctx,
		"", //exchange
		r.queueProducer.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(message),
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("message send to queue %s\n", QueueName)

	return nil
}

func (r *Rabbitmq) Consumer() error {

	q, err := r.ch.QueueDeclare(
		QueueName,
		false,
		false,
		false,
		false,
		nil,
	)

	msgs, err := r.ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	go func() {
		for delivery := range msgs {
			log.Printf("received a message: %s\n", delivery.Body)
			delivery.Ack(false)
		}
	}()

	return nil
}

func (r *Rabbitmq) ReadReturn() {
	for r := range r.notifyReturnCh {
		fmt.Printf("message returned from server: %s\n", r.AppId)
	}
}

func (r *Rabbitmq) ReadConfirmations() {
	for c := range r.confirmationCh {
		fmt.Printf("message confirmed from server. tag: %v, ack: %v\n", c.DeliveryTag, c.Ack)
	}
}

func (r *Rabbitmq) Close() error {

	if err := r.conn.Close(); err != nil {
		return fmt.Errorf("conn.Close: %w", err)
	}

	return nil
}
