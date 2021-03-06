package descriptor

import (
	"google.golang.org/grpc"
	"testing"
)

func Test_descriptor_empty_package(t *testing.T) {
	desc, err := NewDescriptor("../testdata/protobuf/service/service.pb")
	if err != nil {
		t.Fatal(err)
	}

	sds := ServiceDescs(desc)

	service := "Service"
	if s := sds[0].ServiceName; s != service {
		t.Errorf("%s != %s", s, service)
	}

	metadata := "testdata/protobuf/service/service.proto"
	if s := sds[0].Metadata; s != metadata {
		t.Errorf("%s != %s", s, metadata)
	}

	method := "Call"
	if s := sds[0].Methods[0].MethodName; s != method {
		t.Errorf("%s != %s", s, method)
	}
}

func Test_descriptor_ping(t *testing.T) {
	desc, err := NewDescriptor("../testdata/protobuf/ping/ping.pb")
	if err != nil {
		t.Fatal(err)
	}

	sds := ServiceDescs(desc)

	service := "ping.PingService"
	if s := sds[0].ServiceName; s != service {
		t.Errorf("%s != %s", s, service)
	}

	metadata := "testdata/protobuf/ping/ping.proto"
	if s := sds[0].Metadata; s != metadata {
		t.Errorf("%s != %s", s, metadata)
	}

	method := "Send"
	if s := sds[0].Methods[0].MethodName; s != method {
		t.Errorf("%s != %s", s, method)
	}

	checkStream := func(sd grpc.StreamDesc, name string, c bool, s bool) {
		if s := sd.StreamName; s != name {
			t.Errorf("%s != %s %v", s, name, sd)
		}
		if b := sd.ClientStreams; b != c {
			t.Errorf("%v != %v %v", b, c, sd)
		}
		if b := sd.ServerStreams; b != s {
			t.Errorf("%v != %v %v", b, s, sd)
		}
	}

	checkStream(sds[0].Streams[0], "SendStreamC", true, false)
	checkStream(sds[0].Streams[1], "SendStreamS", false, true)
	checkStream(sds[0].Streams[2], "SendStreamB", true, true)
}
