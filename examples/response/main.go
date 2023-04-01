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
func (h *HelloWorld) Receive(e *actor.Environ, p *actor.Parcel) {
	switch msg := p.Message.(type) {
	case *Message:
		fmt.Printf("receive: %s\n", msg.Text)

		fmt.Println("set response value")
		p.Response.SetValue(fmt.Sprintf("Hello from '%s'", e.Self()))
	}
}

func main() {
	// 1. create engine
	engine := actor.NewEngine()

	// 2. spawn the actor
	actorId, err := engine.Spawn(&HelloWorld{}, "actor_1")
	if err != nil {
		fmt.Printf("unable to spawn the actor: %s\n", err)
		os.Exit(1)
	}

	// 3. send a messages with response
	fmt.Println("send message")
	result := engine.SendWithResponse(actorId, &Message{Text: "Hello actor!"})
	if result == nil {
		fmt.Println("message was not sent")
	}

	// 4. read response value
	value, err := result.Result(context.Background())
	if err != nil {
		fmt.Printf("failed to read response value: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("response value: %s\n", value)

	// press Ctrl+C to exit
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)
	<-exit

	// 5. gracefull shutdown
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := engine.Shutdown(ctx); err != nil {
		fmt.Println(err)
	}
}
