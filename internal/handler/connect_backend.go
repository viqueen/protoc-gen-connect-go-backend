package handler

import (
	"errors"
	"fmt"
	"github.com/viqueen/go-protoc-gen-plugin/internal/codegen"
	"github.com/viqueen/go-protoc-gen-plugin/internal/helpers"
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

	apiTarget := toApiTarget(packageName)
	for _, service := range services {
		serviceName := service.GetName()
		serviceNameSnake := helpers.CamelToSnake(serviceName)
		serviceFileName := fmt.Sprintf("%s.go", serviceNameSnake)
		serviceFilePath := filepath.Join("internal", apiTarget, serviceFileName)

		connectServiceFileContent, serviceErr := codegen.ConnectServiceFile(codegen.ConnectServiceFileInput{
			PackageName:    packageName,
			ApiPackage:     apiPackage,
			DataGenPackage: dataGenPackage,
		}, service)
		if serviceErr != nil {
			response.Error = proto.String(serviceErr.Error())
		}

		response.File = append(response.File, &pluginpb.CodeGeneratorResponse_File{
			Name:    proto.String(serviceFilePath),
			Content: proto.String(connectServiceFileContent),
		})

		for _, method := range service.GetMethod() {
			rpcName := method.GetName()
			rpcFileName := fmt.Sprintf("%s.go", helpers.CamelToSnake(rpcName))
			rpcFilePath := filepath.Join("internal", apiTarget, rpcFileName)

			rpcFileContent, rpcErr := codegen.ConnectHandlerFile(codegen.ConnectHandlerFileInput{
				PackageName: packageName,
				ApiPackage:  apiPackage,
			}, service, method)
			if rpcErr != nil {
				response.Error = proto.String(rpcErr.Error())
			}

			response.File = append(response.File, &pluginpb.CodeGeneratorResponse_File{
				Name:    proto.String(rpcFilePath),
				Content: proto.String(rpcFileContent),
			})
		}
	}

	return nil
}
