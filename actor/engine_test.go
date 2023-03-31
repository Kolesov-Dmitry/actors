package actor

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	__responseMessage = "response"
)

type __testMessage struct {
	done     chan struct{}
	response string
}

type __testReceiver struct{}

func (*__testReceiver) Receive(_ *Environ, p *Parcel) {
	msg, ok := p.Message.(*__testMessage)
	if ok {
		close(msg.done)
		if p.Response != nil {
			p.Response.SetValue(msg.response)
		}
	}
}

func Test_Spawn(t *testing.T) {
	engine := NewEngine()

	expectedActors := 10
	actorsIds := make([]ID, expectedActors)

	t.Run("Spawn", func(t *testing.T) {
		for idx := 0; idx < expectedActors; idx++ {
			id, err := engine.Spawn(&__testReceiver{}, "actor", "test", strconv.Itoa(idx))
			assert.Nil(t, err)
			assert.False(t, id.IsEmpty())

			actorsIds[idx] = id
		}

		assert.Equal(t, expectedActors, len(engine.disp.actors))
	})

	t.Run("Spawn with empty receiver", func(t *testing.T) {
		id, err := engine.Spawn(nil, "empty_receiver")
		assert.NotNil(t, err)
		assert.True(t, id.IsEmpty())

		id, err = engine.Spawn(&__testReceiver{}, "")
		assert.NotNil(t, err)
		assert.True(t, id.IsEmpty())
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

	id, err := engine.Spawn(&__testReceiver{}, "test")
	require.Nil(t, err)
	require.NotNil(t, id)

	ok := engine.Send(id, &__testMessage{
		done:     done,
		response: __responseMessage,
	})
	assert.True(t, ok)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		t.Fatal("message wasn't sent")
	case <-done:
	}

	err = engine.Drop(ctx, id)
	assert.Nil(t, err)

	err = engine.Shutdown(ctx)
	assert.Nil(t, err)
}

func Test_SendWithResponse(t *testing.T) {
	engine := NewEngine()

	done := make(chan struct{})

	id, err := engine.Spawn(&__testReceiver{}, "test")
	require.Nil(t, err)
	require.NotNil(t, id)

	response := engine.SendWithResponse(id, &__testMessage{
		done:     done,
		response: __responseMessage,
	})
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

const __middlewareResponse = "middleware"

func __testMiddleware(next ReceiveFunc) ReceiveFunc {
	return func(env *Environ, p *Parcel) {
		if msg, ok := p.Message.(*__testMessage); ok {
			msg.response = __middlewareResponse
		}

		next(env, p)
	}
}

func Test_Middleware(t *testing.T) {
	engine := NewEngine(WithMiddleware(
		__testMiddleware,
	))

	done := make(chan struct{})

	id, err := engine.Spawn(&__testReceiver{}, "test")
	require.Nil(t, err)
	require.NotNil(t, id)

	response := engine.SendWithResponse(id, &__testMessage{
		done:     done,
		response: __responseMessage,
	})
	require.NotNil(t, response)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	value, err := response.Result(ctx)
	require.Nil(t, err)

	responseMessage, ok := value.(string)
	require.True(t, ok)
	assert.Equal(t, __middlewareResponse, responseMessage)

	err = engine.Drop(ctx, id)
	assert.Nil(t, err)

	err = engine.Shutdown(ctx)
	assert.Nil(t, err)
}

type __broadcastMessage struct {
	doneOne chan struct{}
	doneTwo chan struct{}
}

type __broadcastOneReceiver struct{}

func (*__broadcastOneReceiver) Receive(_ *Environ, p *Parcel) {
	msg, ok := p.Message.(*__broadcastMessage)
	if ok {
		close(msg.doneOne)
	}
}

type __broadcastTwoReceiver struct{}

func (*__broadcastTwoReceiver) Receive(_ *Environ, p *Parcel) {
	msg, ok := p.Message.(*__broadcastMessage)
	if ok {
		close(msg.doneTwo)
	}
}

func Test_Broadcast(t *testing.T) {
	engine := NewEngine()

	doneOne := make(chan struct{})
	doneTwo := make(chan struct{})

	idone, err := engine.Spawn(&__broadcastOneReceiver{}, "test", "one")
	require.Nil(t, err)
	require.NotNil(t, idone)

	idtwo, err := engine.Spawn(&__broadcastTwoReceiver{}, "test", "two")
	require.Nil(t, err)
	require.NotNil(t, idtwo)

	group := NewBroadcastGroup(idone, idtwo)

	ok := engine.Broadcast(group,
		&__broadcastMessage{
			doneOne: doneOne,
			doneTwo: doneTwo,
		},
	)
	require.True(t, ok)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		t.Fatal("message wasn't sent to reciever one")
	case <-doneOne:
	}

	select {
	case <-ctx.Done():
		t.Fatal("message wasn't sent sent to reciever two")
	case <-doneTwo:
	}

	err = engine.Shutdown(ctx)
	assert.Nil(t, err)
}
