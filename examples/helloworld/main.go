package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/klsvdm/actors/actor"
)

// Define type of the message
// which will be sent to the actor
type Message struct {
	Text string
}

// Define the actor type
type HelloWorld struct{}

// Actor should implement Receive method
func (h *HelloWorld) Receive(_ *actor.Environ, p *actor.Parcel) {
	switch msg := p.Message.(type) {
	case *Message:
		fmt.Printf("receive: %s\n", msg.Text)
	}
}

func main() {
	// 1. create engine
	engine := actor.NewEngine()

	// 2. spawn the actor
	actorId, err := engine.Spawn(&HelloWorld{}, "hello_world")
	if err != nil {
		fmt.Printf("unable to spawn the actor: %s\n", err)
		os.Exit(1)
	}

	// 3. send a bunch of messages
	for i := 1; i <= 10; i++ {
		message := &Message{
			Text: fmt.Sprintf("Hello World! - %d", i),
		}

		fmt.Printf("send: %s\n", message.Text)
		if ok := engine.Send(actorId, message); !ok {
			fmt.Println("message was not sent")
		}
	}

	// press Ctrl+C to exit
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)
	<-exit

	// 4. gracefull shutdown
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := engine.Shutdown(ctx); err != nil {
		fmt.Println(err)
	}
}
