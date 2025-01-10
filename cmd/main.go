package main

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
	"io"
	"log"
	"os"
)

func main() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("failed to read from stdin: %v", err)
	}
	request := &pluginpb.CodeGeneratorRequest{}
	if err = proto.Unmarshal(data, request); err != nil {
		log.Fatalf("failed to unmarshal input: %v", err)
	}

	response := &pluginpb.CodeGeneratorResponse{}
	for _, protoFile := range request.GetProtoFile() {
		processFile(protoFile, response)
	}
	respond(response)
}

func processFile(protoFile *descriptorpb.FileDescriptorProto, response *pluginpb.CodeGeneratorResponse) {
	services := protoFile.GetService()
	for _, service := range services {
		log.Printf("service: %v", service)
		methods := service.GetMethod()
		for _, method := range methods {
			log.Printf("method: %v", method)
		}
	}
}

func respond(resp *pluginpb.CodeGeneratorResponse) {
	out, err := proto.Marshal(resp)
	if err != nil {
		log.Fatalf("Failed to marshal response: %v", err)
	}
	_, err = os.Stdout.Write(out)
	if err != nil {
		log.Fatalf("Failed to write response: %v", err)
	}
}
