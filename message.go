package main

type message struct {
	bytes []byte
}

func (*message) ProtoMessage() {}

func (m *message) Reset() {
	*m = message{}
}

func (m *message) String() string {
	return ""
}
