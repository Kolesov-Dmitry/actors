package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Kolesov-Dmitry/actors/actor"
)

type Foo struct {
}

func (*Foo) Receive(_ *actor.Environ, p *actor.Parcel) {
	msg, ok := p.Message().(string)
	if ok {
		fmt.Println(msg)
		p.Respond("Hi")
	}
}

const (
	produrersSize = 5
)

func main() {
	e := actor.NewEngine(actor.WithCapacity(5000))
	id := e.Spawn(&Foo{}, "test")

	resp := e.SendWithResponse(id, "Say Hi!")

	value, err := resp.Result(context.Background())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(value)

	wg := &sync.WaitGroup{}
	wg.Add(produrersSize)

	for idx := 0; idx < produrersSize; idx++ {
		go func() {
			defer wg.Done()

			msg := "Hello"
			for idx := 0; idx < 100_000; idx++ {
				e.Send(id, msg)
			}
		}()
	}

	start := time.Now()
	wg.Wait()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		fmt.Println(err)
	}

	elapsed := time.Since(start)
	fmt.Printf("Elapsed: %s\n", elapsed)
}
