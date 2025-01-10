package handler

import (
	"errors"
	"fmt"
	"github.com/viqueen/go-protoc-gen-plugin/internal/codegen"
	"github.com/viqueen/protoc-gen-sqlc/pkg/helpers"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
	"path/filepath"
	"strings"
)

func requestMappers(params map[string]string, protoFile *descriptorpb.FileDescriptorProto, response *pluginpb.CodeGeneratorResponse) error {
	messages := protoFile.GetMessageType()
	if len(messages) == 0 {
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

	packageName := protoFile.GetPackage()
	apiTarget := toApiTarget(packageName)

	for _, message := range messages {
		_, sqlRequestOk := helpers.SqlcRequestOption(message)
		if !sqlRequestOk {
			continue
		}
		requestMappersFileName := fmt.Sprintf("request_mapper_%s.go", strings.ToLower(message.GetName()))
		requestMappersFilePath := filepath.Join("internal", apiTarget, requestMappersFileName)
		requestMappersFileContent, err := codegen.RequestMapperFile(codegen.RequestMapperFileInput{
			PackageName:    packageName,
			ApiPackage:     apiPackage,
			DataGenPackage: dataGenPackage,
		}, message)
		if err != nil {
			response.Error = proto.String(err.Error())
		}
		response.File = append(response.File, &pluginpb.CodeGeneratorResponse_File{
			Name:    proto.String(requestMappersFilePath),
			Content: proto.String(requestMappersFileContent),
		})
	}
	return nil
}
