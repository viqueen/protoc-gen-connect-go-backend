package main

import (
	"github.com/viqueen/protoc-gen-connect-go-backend/protoc-gen-connect-go-backend/handler"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
	"io"
	"log"
	"os"
	"strings"
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
	params := asMap(request.GetParameter())
	response := &pluginpb.CodeGeneratorResponse{}
	for _, protoFile := range request.GetProtoFile() {
		err = handler.ProtoFileHandler(params, protoFile, response)
		if err != nil {
			response.Error = proto.String(err.Error())
		}
	}
	respond(response)
}

func asMap(params string) map[string]string {
	result := make(map[string]string)
	for _, param := range strings.Split(params, ",") {
		parts := strings.Split(param, "=")
		if len(parts) != 2 {
			continue
		}
		result[parts[0]] = parts[1]
	}
	return result
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
