package actor

type Parcel struct {
	engine   *Engine
	response *Response
	msg      any
}

func newParcel(engine *Engine, response *Response, msg any) *Parcel {
	return &Parcel{
		engine:   engine,
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
