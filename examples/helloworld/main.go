package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/klsvdm/actors/actor"
)

type Message struct {
	Text string
}

type HelloWorld struct{}

func (h *HelloWorld) Receive(_ *actor.Environ, p *actor.Parcel) {
	switch msg := p.Message.(type) {
	case *Message:
		fmt.Printf("receive: %s\n", msg.Text)
	}
}

func main() {
	engine := actor.NewEngine()

	// 1. spawn the actor
	actorId, err := engine.Spawn(&HelloWorld{}, "hello_world")
	if err != nil {
		fmt.Printf("unable to spawn the actor: %s\n", err)
		os.Exit(1)
	}

	// 2. send a bunch of messages
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
}
