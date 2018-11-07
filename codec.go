package main

type codec struct{}

func (codec) Marshal(v interface{}) ([]byte, error) {
	return v.(*message).bytes, nil
}

func (codec) Unmarshal(data []byte, v interface{}) error {
	v.(*message).bytes = data
	return nil
}

func (codec) String() string {
	return "grpc-proxy message codec"
}
