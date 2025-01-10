package handler

import (
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
	"log"
)

func ProtoFileHandler(params map[string]string, protoFile *descriptorpb.FileDescriptorProto, response *pluginpb.CodeGeneratorResponse) error {
	dataMappersErr := sqlcDataMappers(params, protoFile, response)
	if dataMappersErr != nil {
		log.Fatalf("failed to generate data mappers: %v", dataMappersErr)
	}
	requestMappersErr := requestMappers(params, protoFile, response)
	if requestMappersErr != nil {
		log.Fatalf("failed to generate request mappers: %v", requestMappersErr)
	}
	connectBackendErr := connectBackend(params, protoFile, response)
	if connectBackendErr != nil {
		log.Fatalf("failed to generate connect backend: %v", connectBackendErr)
	}
	return nil
}
