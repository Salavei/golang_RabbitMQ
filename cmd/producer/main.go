package main

import (
	"context"
	"github.com/Salavei/golang_RabbitMQ/internal"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

func main() {
	conn, err := internal.ConnectRabbitMQ(
		"salavei", "root123",
		"localhost:5672", "customers")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client, err := internal.NewRabbitMQClient(conn)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	err = client.CreateQueue("customers_created", true, false)
	if err != nil {
		panic(err)
	}
	// If true == false. The data will not be saved
	err = client.CreateQueue("customers_test", false, true)
	if err != nil {
		panic(err)
	}

	err = client.CreateBinding("customers_created", "customers.created.*", "customer_events")
	if err != nil {
		panic(err)
	}

	err = client.CreateBinding("customers_test", "customers.*", "customer_events")
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	for i := 0; i < 10; i++ {
		err = client.Send(ctx, "customer_events", "customers.created.us", amqp.Publishing{
			ContentType:     "text/plain",
			DeliveryMode:    amqp.Persistent,
			Body:            []byte(`An cool message between services`),
		})
		if err != nil {
			panic(err)
		}
		// Sending a transient message
		err = client.Send(ctx, "customer_events", "customers.test", amqp.Publishing{
			ContentType:     "text/plain",
			DeliveryMode:    amqp.Transient,
			Body:            []byte(`An uncool undurable message`),
		})
		if err != nil {
			panic(err)
		}
	}

	log.Println(client)
}
