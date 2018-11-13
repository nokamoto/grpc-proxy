protobuf = testdata/protobuf

ping_proto = $(protobuf)/ping/ping.proto
ping_pb = $(protobuf)/ping/ping.pb
ping_pb_go = test/ping.pb.go

service_proto = $(protobuf)/service/service.proto
service_pb = $(protobuf)/service/service.pb

objs = $(ping_pb) $(ping_pb_go) $(service_pb)

all: $(objs)

$(ping_pb): $(ping_proto)
	prototool format -d $(ping_proto) || prototool format -w $(ping_proto)
	protoc --include_imports --include_source_info $(ping_proto) --descriptor_set_out $(ping_pb)

$(ping_pb_go): $(ping_proto)
	protoc --go_out=plugins=grpc:${GOPATH}/src $(ping_proto)

$(service_pb): $(service_proto)
	prototool format -d $(service_proto) || prototool format -w $(service_proto)
	protoc --include_imports --include_source_info $(service_proto) --descriptor_set_out $(service_pb)

clean:
	rm $(objs)

test: all
	dep check
	go test ./...
	test -z `go fmt ./...`
	golint -set_exit_status `go list ./... | grep -v /vendor/`
