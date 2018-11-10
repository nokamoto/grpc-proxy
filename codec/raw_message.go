package codec

type RawMessage struct {
	bytes []byte
}

func (*RawMessage) ProtoMessage() {}

func (m *RawMessage) Reset() {
	*m = RawMessage{}
}

func (*RawMessage) String() string {
	return ""
}
