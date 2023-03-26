package actor

type Parcel struct {
	engine   *Engine
	sender   *ID
	response *Response
	msg      any
}

func newParcel(engine *Engine, sender *ID, response *Response, msg any) *Parcel {
	return &Parcel{
		engine:   engine,
		sender:   sender,
		response: response,
		msg:      msg,
	}
}

func (p *Parcel) Message() any {
	return p.msg
}

func (p *Parcel) Respond(value any) {
	if p.response == nil {
		return
	}

	p.response.setValue(value)
}
