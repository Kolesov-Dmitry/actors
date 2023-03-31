package actor

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type __testReceiverWithSend struct{}

func (*__testReceiverWithSend) Receive(e *Environ, p *Parcel) {
	if !p.Sender.IsEmpty() {
		e.Send(p.Sender, p.Message)
	}
}

func Test_EnvironSpawnAndDropChild(t *testing.T) {
	engine := NewEngine()
	parentId, err := engine.Spawn(&__testReceiver{}, "parent")
	require.Nil(t, err)
	require.NotEqual(t, uuid.UUID{}, parentId)

	parent := engine.disp.ActorById(parentId).(*actor)
	require.NotNil(t, parent)

	var childId ID

	t.Run("Spawn", func(t *testing.T) {
		childId, err = parent.environ.SpawnChild(&__testReceiver{}, "child")
		require.Nil(t, err)
		require.NotEqual(t, uuid.UUID{}, childId)

		assert.Len(t, parent.children, 1)

		_, ok := parent.children[childId]
		assert.True(t, ok)
	})

	t.Run("Drop", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		err := parent.environ.DropChild(ctx, childId)
		require.Nil(t, err)

		assert.Len(t, parent.children, 0)
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err = engine.Shutdown(ctx)
	assert.Nil(t, err)
}

func Test_EnvironSend(t *testing.T) {
	engine := NewEngine()
	receiverId, err := engine.Spawn(&__testReceiver{}, "receiver")
	require.Nil(t, err)
	require.NotNil(t, receiverId)

	receiver := engine.disp.ActorById(receiverId).(*actor)
	require.NotNil(t, receiver)

	senderId, err := engine.Spawn(&__testReceiverWithSend{}, "sender")
	require.Nil(t, err)
	require.NotNil(t, senderId)

	done := make(chan struct{})

	ok := receiver.environ.Send(senderId, &__testMessage{done: done})
	assert.True(t, ok)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		t.Fatal("message wasn't sent")
	case <-done:
	}

	err = engine.Shutdown(ctx)
	assert.Nil(t, err)
}
