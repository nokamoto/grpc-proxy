package codec

import (
	"encoding/hex"
)

// RawMessage represents arbitrary gRPC messages.
type RawMessage struct {
	bytes []byte
}

// ProtoMessage is nothing to do.
func (*RawMessage) ProtoMessage() {}

// Reset resets the bytes.
func (m *RawMessage) Reset() {
	*m = RawMessage{}
}

func (m *RawMessage) String() string {
	return hex.EncodeToString(m.bytes)
}
