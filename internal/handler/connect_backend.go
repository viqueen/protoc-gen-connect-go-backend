package handler

import (
	"errors"
	"github.com/viqueen/go-protoc-gen-plugin/internal/codegen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
	"path/filepath"
)

func connectBackend(params map[string]string, protoFile *descriptorpb.FileDescriptorProto, response *pluginpb.CodeGeneratorResponse) error {
	services := protoFile.GetService()
	if len(services) == 0 {
		return nil
	}

	apiPackage, ok := params["api_package"]
	if !ok {
		return errors.New("api_package is required")
	}
	dataGenPackage, ok := params["data_gen_package"]
	if !ok {
		return errors.New("data_gen_package is required")
	}
	internalPackage, ok := params["internal_package"]
	if !ok {
		return errors.New("internal_package is required")
	}

	packageName := protoFile.GetPackage()

	serverFileName := "main.go"
	serverFilePath := filepath.Join("cmd", "server", serverFileName)
	connectServerFileContent, err := codegen.ConnectServerFile(codegen.ConnectServerFileInput{
		DataGenPackage:  dataGenPackage,
		PackageName:     packageName,
		InternalPackage: internalPackage,
		ApiPackage:      apiPackage,
	}, services)
	if err != nil {
		response.Error = proto.String(err.Error())
	}
	response.File = append(response.File, &pluginpb.CodeGeneratorResponse_File{
		Name:    proto.String(serverFilePath),
		Content: proto.String(connectServerFileContent),
	})

	return nil
}
