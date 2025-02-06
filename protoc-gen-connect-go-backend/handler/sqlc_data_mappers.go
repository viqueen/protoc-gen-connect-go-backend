package handler

import (
	"errors"
	"fmt"
	"github.com/viqueen/protoc-gen-connect-go-backend/protoc-gen-connect-go-backend/codegen"
	"github.com/viqueen/protoc-gen-sqlc/pkg/helpers"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
	"path/filepath"
	"strings"
)

func sqlcDataMappers(params map[string]string, protoFile *descriptorpb.FileDescriptorProto, response *pluginpb.CodeGeneratorResponse) error {
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
	dataGenTablePrefix, _ := params["data_gen_table_prefix"]

	packageName := protoFile.GetPackage()
	apiTarget := toApiTarget(protoFile.GetPackage())
	for _, message := range messages {
		_, sqlEntityOk := helpers.SqlcEntityOption(message)
		if !sqlEntityOk {
			continue
		}
		dataMapperFileName := fmt.Sprintf("data_mapper_%s.go", strings.ToLower(message.GetName()))
		dataMapperFilePath := filepath.Join("internal", apiTarget, dataMapperFileName)
		dataMapperFileContent, err := codegen.DataMapperFile(codegen.DataMapperFileInput{
			PackageName:     packageName,
			ApiPackage:      apiPackage,
			DataGenPackage:  dataGenPackage,
			TableNamePrefix: dataGenTablePrefix,
		}, message)
		if err != nil {
			response.Error = proto.String(err.Error())
		}
		response.File = append(response.File, &pluginpb.CodeGeneratorResponse_File{
			Name:    proto.String(dataMapperFilePath),
			Content: proto.String(dataMapperFileContent),
		})
	}

	return nil
}
