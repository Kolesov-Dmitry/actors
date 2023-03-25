package actor

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Subscribe(t *testing.T) {
	stream := newEventStream(1)
	stream.Start()

	expectedSubs := 10
	subsIds := make([]uuid.UUID, expectedSubs)

	t.Run("Subscribe", func(t *testing.T) {
		for idx := 0; idx < expectedSubs; idx++ {
			id := stream.Subscribe(func(events []any) {})
			assert.NotEqual(t, uuid.UUID{}, id)

			subsIds[idx] = id
		}

		assert.Equal(t, expectedSubs, len(stream.subs))
	})

	t.Run("Subscribe without event handler", func(t *testing.T) {
		id := stream.Subscribe(nil)
		assert.Equal(t, uuid.UUID{}, id)
		assert.Equal(t, expectedSubs, len(stream.subs))
	})

	t.Run("Unsubscribe", func(t *testing.T) {
		for _, id := range subsIds {
			stream.Unsubscribe(id)
		}

		assert.Equal(t, 0, len(stream.subs))
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := stream.Stop(ctx)
	assert.Nil(t, err)
}

func Test_Dispatch(t *testing.T) {
	stream := newEventStream(1)
	stream.Start()

	message := "Test Test"

	done := make(chan struct{})

	id := stream.Subscribe(func(events []any) {
		if len(events) == 0 {
			return
		}

		require.Len(t, events, 1)
		text, ok := events[0].(string)
		require.True(t, ok)
		assert.Equal(t, message, text)

		close(done)
	})

	stream.Dispatch(message)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		t.Log("message wasn't dispatched")
		t.Fail()
	case <-done:
	}

	stream.Unsubscribe(id)
	err := stream.Stop(ctx)
	assert.Nil(t, err)
}
