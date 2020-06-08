package main

import (
	"context"
	_ "fmt"
	"log"
	"time"

	"github.com/Azure/go-amqp"
)

func main() {
	// Create client
	log.Printf("Dialing AMQP server: %s", "amqps://my-namespace.servicebus.windows.net")
	client, err := amqp.Dial("amqps://my-namespace.servicebus.windows.net",
		amqp.ConnSASLPlain("access-key-name", "access-key"),
	)
	if err != nil {
		log.Fatal("Dialing AMQP server:", err)
	}
	defer client.Close()

	// Open a session
	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Creating AMQP session:", err)
	}

	ctx := context.Background()

	// Send a message
	{
		// Create a sender
		sender, err := session.NewSender(
			amqp.LinkTargetAddress("/queue-name"),
		)
		if err != nil {
			log.Fatal("Creating sender link:", err)
		}

		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

		// Send message
		err = sender.Send(ctx, amqp.NewMessage([]byte("Hello!")))
		if err != nil {
			log.Fatal("Sending message:", err)
		}

		sender.Close(ctx)
		cancel()
	}
}