package main

import (
	"context"
	"github.com/Salavei/golang_RabbitMQ/internal"
	"golang.org/x/sync/errgroup"
	"log"
	"time"
)

func main() {
	conn, err := internal.ConnectRabbitMQ("salavei", "root123",
		"localhost:5672", "customers")
	defer conn.Close()
	if err != nil {
		panic(err)
	}

	client, err := internal.NewRabbitMQClient(conn)
	defer client.Close()
	if err != nil {
		panic(err)
	}

	messageBus, err := client.Consume("customers_created", "email-service", true)
	if err != nil {
		panic(err)
	}
	var blocking chan struct{}
	// Set a timeout for 15secs.
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 15 * time.Second)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	// errgroup allows us concurrent task.
	g.SetLimit(10)

	go func() {
		for message := range messageBus {
			// Spawn a worker.
			msg := message
			g.Go(func() error {
				log.Printf("NewMessage: %v", msg)
				time.Sleep(10 * time.Second)
				err = msg.Ack(false)
				if err != nil {
					log.Println("Ack message failed")
					return err
				}
				log.Printf("Ackowledged message %s\n", message.MessageId)
				return nil
			})
		}
	}()
	log.Println("Consuming, use CTRL+C to exit")
	<-blocking
}
