package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Kolesov-Dmitry/actors/actor"
)

type Foo struct {
}

func (*Foo) Receive(p *actor.Parcel) {
	msg, ok := p.Message().(string)
	if ok {
		fmt.Println(msg)
	}
}

func main() {
	e := actor.NewEngine(actor.WithCapacity(5000))
	id := e.Spawn(&Foo{}, "test")

	wg := &sync.WaitGroup{}
	wg.Add(10)

	for idx := 0; idx < 10; idx++ {
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
