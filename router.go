package main

import (
	pb "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

type router struct {
	clusters map[string]*cluster
}

func newRouter(fds *pb.FileDescriptorSet, routes *yamlRoutes, clusters *yamlClusters) (*router, error) {
	for _, fd := range fds.File {
		fd.GetService
	}
}
