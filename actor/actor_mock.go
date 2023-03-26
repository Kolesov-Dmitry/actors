package actor

import "context"

type actorMock struct {
	id *ID
}

func (a *actorMock) ID() *ID {
	return a.id
}

func (a *actorMock) Invoke(_ *Parcel) {
}

func (a *actorMock) AddChild(_ Actor) {
}

func (a *actorMock) DropChild(_ context.Context, _ *ID) error {
	return nil
}

func (a *actorMock) Shutdown(_ context.Context) error {
	return nil
}
