package actor

import "strings"

type ID struct {
	addr string
	id   string
}

const (
	addressTagsSeparator = "/"
)

func newID(addr, id string, tags ...string) *ID {
	actorId := &ID{
		addr: addr,
		id:   id,
	}

	if len(tags) != 0 {
		actorId.id += addressTagsSeparator + strings.Join(tags, addressTagsSeparator)
	}

	return actorId
}

func (id *ID) String() string {
	return id.addr + addressTagsSeparator + id.id
}

func (a *ID) Equals(other *ID) bool {
	return a.addr == other.addr && a.id == other.id
}
