package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/klsvdm/actors/actor"
)

// Define the actor type
type Actor struct{}

// Actor should implement Receive method
func (a *Actor) Receive(e *actor.Environ, p *actor.Parcel) {
	switch p.Message.(type) {
	case actor.StartedEvent:
		fmt.Printf("actor '%s' is started\n", e.Self())

	case actor.AboutToStopEvent:
		fmt.Printf("actor '%s' is about to stop\n", e.Self())
	}
}

func main() {
	// 1. create engine
	engine := actor.NewEngine()

	// 2. spawn the actor
	fmt.Println("spawn the actor")

	actorId, err := engine.Spawn(&Actor{}, "actor-with-hooks")
	if err != nil {
		fmt.Printf("unable to spawn the actor: %s\n", err)
		os.Exit(1)
	}

	// 3. wait two seconds
	time.Sleep(2 * time.Second)

	// 4. drop the actor
	fmt.Println("drop the actor")

	if err := engine.Drop(context.Background(), actorId); err != nil {
		fmt.Printf("unable to drop the actor: %s\n", err)
		os.Exit(1)
	}

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
