package actor

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	__responseMessage = "response"
)

type __testMessage struct {
	done chan struct{}
}

type __testReceiver struct{}

func (*__testReceiver) Receive(_ *Environ, p *Parcel) {
	msg, ok := p.Message.(*__testMessage)
	if ok {
		close(msg.done)
		if p.Response != nil {
			p.Response.SetValue(__responseMessage)
		}
	}
}

func Test_Spawn(t *testing.T) {
	engine := NewEngine()

	expectedActors := 10
	actorsIds := make([]*ID, expectedActors)

	t.Run("Spawn", func(t *testing.T) {
		for idx := 0; idx < expectedActors; idx++ {
			id := engine.Spawn(&__testReceiver{}, fmt.Sprintf("actor/test/%d", idx))
			assert.NotNil(t, id)

			actorsIds[idx] = id
		}

		assert.Equal(t, expectedActors, len(engine.disp.actors))
	})

	t.Run("Spawn with empty receiver", func(t *testing.T) {
		id := engine.Spawn(nil, "empty_receiver")
		assert.Nil(t, id)

		id = engine.Spawn(&__testReceiver{}, "")
		assert.Nil(t, id)
	})

	t.Run("Drop", func(t *testing.T) {
		for _, id := range actorsIds {
			err := engine.Drop(context.Background(), id)
			assert.Nil(t, err)
		}

		assert.Equal(t, 0, len(engine.disp.actors))
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := engine.Shutdown(ctx)
	assert.Nil(t, err)
}

func Test_Send(t *testing.T) {
	engine := NewEngine()

	done := make(chan struct{})

	id := engine.Spawn(&__testReceiver{}, "test")
	require.NotNil(t, id)

	ok := engine.Send(id, &__testMessage{done})
	assert.True(t, ok)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		t.Fatal("message wasn't sent")
	case <-done:
	}

	err := engine.Drop(ctx, id)
	assert.Nil(t, err)

	err = engine.Shutdown(ctx)
	assert.Nil(t, err)
}

func Test_SendWithResponse(t *testing.T) {
	engine := NewEngine()

	done := make(chan struct{})

	id := engine.Spawn(&__testReceiver{}, "test")
	require.NotNil(t, id)

	response := engine.SendWithResponse(id, &__testMessage{done})
	require.NotNil(t, response)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	value, err := response.Result(ctx)
	require.Nil(t, err)

	responseMessage, ok := value.(string)
	require.True(t, ok)
	assert.Equal(t, __responseMessage, responseMessage)

	err = engine.Drop(ctx, id)
	assert.Nil(t, err)

	err = engine.Shutdown(ctx)
	assert.Nil(t, err)
}
