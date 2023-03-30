package actor

type ID struct {
	value string
}

func newID(name string, tags ...string) *ID {
	actorId := &ID{
		value: name,
	}

	for _, tag := range tags {
		actorId.value += "/" + tag
	}

	return actorId
}

func (id *ID) String() string {
	return id.value
}
