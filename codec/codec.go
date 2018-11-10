package codec

// RawCodec implements grpc.CustomCodec for RawMessage.
type RawCodec struct{}

// Marshal returns bytes from RawMessage.
func (RawCodec) Marshal(v interface{}) ([]byte, error) {
	return v.(*RawMessage).bytes, nil
}

// Unmarshal returns RawMessage from the bytes.
func (RawCodec) Unmarshal(data []byte, v interface{}) error {
	v.(*RawMessage).bytes = data
	return nil
}

func (RawCodec) String() string {
	return "grpc-proxy raw codec"
}
