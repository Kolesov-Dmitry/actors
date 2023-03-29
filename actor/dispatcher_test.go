package actor

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_DispatcherAdd(t *testing.T) {
	d := newDispatcher()

	expetedLen := 2

	for idx := 0; idx < expetedLen; idx++ {
		d.Add(&actorMock{newID("test", fmt.Sprintf("actor_%d", idx+1))})
	}

	assert.Equal(t, expetedLen, len(d.actors))
}

func Test_DispatcherRemove(t *testing.T) {
	type Test struct {
		Name        string
		RemoveId    *ID
		ExpectedErr error
	}

	d := newDispatcher()

	tt := struct {
		Actors []Actor
		Tests  []Test
	}{
		Actors: []Actor{
			&actorMock{newID("test", "actor_1")},
			&actorMock{newID("test", "actor_2")},
			&actorMock{newID("test", "actor_3")},
		},
		Tests: []Test{
			{
				Name:        "Remove existing actor",
				RemoveId:    newID("test", "actor_1"),
				ExpectedErr: nil,
			},
			{
				Name:        "Remove non existing actor",
				RemoveId:    newID("test", "actor_5"),
				ExpectedErr: ErrActorDoesNotExists,
			},
			{
				Name:        "Remove already removed actor",
				RemoveId:    newID("test", "actor_1"),
				ExpectedErr: ErrActorDoesNotExists,
			},
		},
	}

	for _, a := range tt.Actors {
		d.Add(a)
	}

	for _, test := range tt.Tests {
		t.Run(test.Name, func(t *testing.T) {
			err := d.Remove(context.Background(), test.RemoveId)
			assert.ErrorIs(t, err, test.ExpectedErr)
		})
	}
}

func Test_DispatcherActorById(t *testing.T) {
	id := newID("test", "actor")

	disp := newDispatcher()
	disp.Add(&actorMock{id})

	actor := disp.ActorById(id)
	require.NotNil(t, actor)
	assert.Equal(t, id, actor.ID())

	err := disp.Remove(context.Background(), id)
	assert.Nil(t, err)

	actor = disp.ActorById(id)
	assert.Nil(t, actor)
}
