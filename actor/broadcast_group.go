package actor

type BroadcastGroup struct {
	ids []*ID
}

func NewBroadcastGroup(ids ...*ID) BroadcastGroup {
	return BroadcastGroup{
		ids: ids,
	}
}
