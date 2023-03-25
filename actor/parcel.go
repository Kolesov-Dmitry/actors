package actor

type Parcel struct {
	engine *Engine
	sender *ID
	msg    any
}

func newParcel(engine *Engine, sender *ID, msg any) *Parcel {
	return &Parcel{
		engine: engine,
		sender: sender,
		msg:    msg,
	}
}

func (p *Parcel) Message() any {
	return p.msg
}

func (p *Parcel) Respond(msg any) {
	if p.sender == nil {
		return
	}

	p.engine.Send(p.sender, msg)
}
