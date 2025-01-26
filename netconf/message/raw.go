package message

type Raw struct {
	RPC
}

// NewRaw produces a raw message that is wrapped by the proper outer.
func NewRaw(data string) *Raw {
	var rpc Raw
	ValidateXML(data, Filter{})
	rpc.Data = data
	rpc.MessageID = uuid()
	return &rpc
}
