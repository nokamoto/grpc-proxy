package codec

type RawCodec struct{}

func (RawCodec) Marshal(v interface{}) ([]byte, error) {
	return v.(*RawMessage).bytes, nil
}

func (RawCodec) Unmarshal(data []byte, v interface{}) error {
	v.(*RawMessage).bytes = data
	return nil
}

func (RawCodec) String() string {
	return "grpc-proxy raw codec"
}
