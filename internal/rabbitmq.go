package internal

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type RabbitClient struct {
	// The connection used by the client.
	conn *amqp.Connection
	// Channel is used to process / Send messages.
	ch *amqp.Channel
}

func ConnectRabbitMQ(username, password, host, vhost string) (*amqp.Connection, error) {
	return amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/%s", username, password, host, vhost))
}

func NewRabbitMQClient(conn *amqp.Connection) (RabbitClient, error) {
	ch, err := conn.Channel()
	if err != nil {
		return RabbitClient{}, err
	}

	err = ch.Confirm(false)
	if err != nil {
		return RabbitClient{}, err
	}
	return RabbitClient{
		conn: conn,
		ch: ch,
	}, nil
}

func (rc RabbitClient) Close() error {
	return rc.ch.Close()
}
// CreateQueue will create a new queue based on given cfg.
func (rc RabbitClient) CreateQueue(queueName string, durable, autodelete bool) error {
	_, err := rc.ch.QueueDeclare(queueName, durable, autodelete, false, false, nil)
	return err
}
// CreateBinding will bind current channel to the given exchange using the routing key provided.
func (rc  RabbitClient) CreateBinding(name, binding, exchange string) error {
	// leaving nowait false, having nowait sett to false will make the return an error if its fails to bind.
	return rc.ch.QueueBind(name, binding, exchange, false, nil)
}
// Send is used to publish payloads onto an exchange with the given routing
func (rc RabbitClient) Send(ctx context.Context, exchange, routingKey string, options amqp.Publishing) error {
	confirmation, err := rc.ch.PublishWithDeferredConfirmWithContext(ctx,
		exchange,
		routingKey,
		// Mandatory is used to determine if an error should be returned upon failure
		true,
		// immediate
		false,
		options)

	if err != nil {
		return err
	}

	log.Println(confirmation.Wait())
	return nil
}
// Consume is used to consume a queue
func (rc RabbitClient) Consume(queue, consumer string, autoAck bool) (<-chan amqp.Delivery, error) {
	return rc.ch.Consume(queue, consumer, autoAck, false, false, false, nil)
}