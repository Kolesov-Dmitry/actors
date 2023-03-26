package actor

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_SpawnAndDropChild(t *testing.T) {
	engine := NewEngine()
	parentId := engine.Spawn(&__testReceiver{}, "parent")
	require.NotEqual(t, uuid.UUID{}, parentId)

	parent := engine.disp.ActorById(parentId).(*actor)
	require.NotNil(t, parent)

	var childId *ID

	t.Run("Spawn", func(t *testing.T) {
		childId = parent.environ.SpawnChild(&__testReceiver{}, "child")
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

	err := engine.Shutdown(ctx)
	assert.Nil(t, err)
}
