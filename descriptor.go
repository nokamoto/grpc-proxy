package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	pb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"google.golang.org/grpc"
	"io/ioutil"
)

type descriptor struct {
	*pb.FileDescriptorSet
}

func newDescriptor(file string) (*descriptor, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	fds := &pb.FileDescriptorSet{}

	if err := proto.Unmarshal(bytes, fds); err != nil {
		return nil, err
	}

	return &descriptor{fds}, nil
}

func (d *descriptor) serviceDescriptors() []*grpc.ServiceDesc {
	descs := make([]*grpc.ServiceDesc, 0)
	for _, file := range d.File {
		descs = append(descs, serviceDescriptorsFromFileDescriptor(file)...)
	}
	return descs
}

func serviceDescriptorsFromFileDescriptor(fd *pb.FileDescriptorProto) []*grpc.ServiceDesc {
	descs := make([]*grpc.ServiceDesc, 0)
	for _, sd := range fd.GetService() {
		descs = append(descs, serviceDescriptor(fd, sd))
	}
	return descs
}

func serviceName(fd *pb.FileDescriptorProto, sd *pb.ServiceDescriptorProto) string {
	name := fd.GetPackage()
	if len(name) > 0 {
		name = name + "."
	}
	name = name + sd.GetName()
	return name
}

func fullMethod(fd *pb.FileDescriptorProto, sd *pb.ServiceDescriptorProto, md *pb.MethodDescriptorProto) string {
	return fmt.Sprintf("/%s/%s", serviceName(fd, sd), md.GetName())
}

func method(fd *pb.FileDescriptorProto, sd *pb.ServiceDescriptorProto, md *pb.MethodDescriptorProto) grpc.MethodDesc {
	return grpc.MethodDesc{
		MethodName: md.GetName(),
		Handler:    unaryProxyHandler(fullMethod(fd, sd, md)),
	}
}

func methods(fd *pb.FileDescriptorProto, sd *pb.ServiceDescriptorProto) []grpc.MethodDesc {
	descs := make([]grpc.MethodDesc, 0)
	for _, md := range sd.GetMethod() {
		if !md.GetClientStreaming() && !md.GetServerStreaming() {
			descs = append(descs, method(fd, sd, md))
		}
	}
	return descs
}

func streamB(fd *pb.FileDescriptorProto, sd *pb.ServiceDescriptorProto, md *pb.MethodDescriptorProto) grpc.StreamDesc {
	return grpc.StreamDesc{
		StreamName:    md.GetName(),
		ClientStreams: true,
		ServerStreams: true,
		Handler:       streamBProxyHandler,
	}
}

func streamC(fd *pb.FileDescriptorProto, sd *pb.ServiceDescriptorProto, md *pb.MethodDescriptorProto) grpc.StreamDesc {
	desc := grpc.StreamDesc{
		StreamName:    md.GetName(),
		ClientStreams: true,
	}

	// todo: &desc may cause unexpected behavior.
	desc.Handler = streamCProxyHandler(fullMethod(fd, sd, md), &desc)

	return desc
}

func streamS(fd *pb.FileDescriptorProto, sd *pb.ServiceDescriptorProto, md *pb.MethodDescriptorProto) grpc.StreamDesc {
	desc := grpc.StreamDesc{
		StreamName:    md.GetName(),
		ServerStreams: true,
	}

	// todo: &desc may cause unexpected behavior.
	desc.Handler = streamSProxyHandler(fullMethod(fd, sd, md), &desc)

	return desc
}

func streams(fd *pb.FileDescriptorProto, sd *pb.ServiceDescriptorProto) []grpc.StreamDesc {
	descs := make([]grpc.StreamDesc, 0)
	for _, md := range sd.GetMethod() {
		cs := md.GetClientStreaming()
		ss := md.GetServerStreaming()
		if cs && ss {
			descs = append(descs, streamB(fd, sd, md))
		} else if cs {
			descs = append(descs, streamC(fd, sd, md))
		} else if ss {
			descs = append(descs, streamS(fd, sd, md))
		}
	}
	return descs
}

func serviceDescriptor(fd *pb.FileDescriptorProto, sd *pb.ServiceDescriptorProto) *grpc.ServiceDesc {
	return &grpc.ServiceDesc{
		ServiceName: serviceName(fd, sd),
		Metadata:    fd.GetName(),
		Methods:     methods(fd, sd),
		Streams:     streams(fd, sd),
		HandlerType: (*proxyServer)(nil),
	}
}
