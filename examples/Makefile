all: example.pb example.go

example.pb: example.proto
	protoc --include_imports --include_source_info example.proto --descriptor_set_out example.pb

example.go: example.proto
	protoc --go_out=plugins=grpc:. example.proto

clean:
	rm example.pb example.go
