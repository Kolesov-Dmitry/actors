package actor

type ID struct {
	addr string
	name string
}

const (
	actorIdSeparator = "/"
)

func newID(addr, name string) *ID {
	actorId := &ID{
		addr: addr,
		name: name,
	}

	return actorId
}

func (id *ID) String() string {
	return id.addr + actorIdSeparator + id.name
}

func (a *ID) Equals(other *ID) bool {
	return a.addr == other.addr && a.name == other.name
}
