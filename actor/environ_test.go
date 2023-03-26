package actor

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_SpawnChild(t *testing.T) {
	engine := NewEngine()
	parentId := engine.Spawn(&__testReceiver{}, "parent")
	require.NotEqual(t, uuid.UUID{}, parentId)

	parent := engine.disp.ActorById(parentId).(*actor)
	require.NotNil(t, parent)

	childId := parent.environ.SpawnChild(&__testReceiver{}, "child")
	require.NotEqual(t, uuid.UUID{}, childId)

	assert.Len(t, parent.children, 1)

	_, ok := parent.children[childId]
	assert.True(t, ok)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := engine.Shutdown(ctx)
	assert.Nil(t, err)
}
