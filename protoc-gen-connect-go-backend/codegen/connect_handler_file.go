package codegen

import (
	"bytes"
	"fmt"
	"github.com/viqueen/protoc-gen-connect-go-backend/protoc-gen-connect-go-backend/helpers"
	"google.golang.org/protobuf/types/descriptorpb"
	"strings"
	"text/template"
)

type ConnectHandlerFileInput struct {
	PackageName string
	ApiPackage  string
}

func ConnectHandlerFile(input ConnectHandlerFileInput, service *descriptorpb.ServiceDescriptorProto, method *descriptorpb.MethodDescriptorProto) (string, error) {
	params := extractConnectHandlerFileParams(input, service, method)
	if params.RpcMethod == "Get" {
		tmpl, err := template.New("getConnectHandlerFile").Parse(getConnectHandlerFileTemplate)
		if err != nil {
			return "", err
		}
		var buf bytes.Buffer
		err = tmpl.Execute(&buf, params)
		if err != nil {
			panic(err)
		}
		return buf.String(), nil
	}
	if params.RpcMethod == "List" {
		tmpl, err := template.New("listConnectHandlerFile").Parse(listConnectHandlerFileTemplate)
		if err != nil {
			return "", err
		}
		var buf bytes.Buffer
		err = tmpl.Execute(&buf, params)
		if err != nil {
			panic(err)
		}
		return buf.String(), nil
	}
	if params.RpcMethod == "Create" {
		tmpl, err := template.New("createConnectHandlerFile").Parse(createConnectHandlerFileTemplate)
		if err != nil {
			return "", err
		}
		var buf bytes.Buffer
		err = tmpl.Execute(&buf, params)
		if err != nil {
			panic(err)
		}
		return buf.String(), nil
	}
	if params.RpcMethod == "Update" {
		tmpl, err := template.New("updateConnectHandlerFile").Parse(updateConnectHandlerFileTemplate)
		if err != nil {
			return "", err
		}
		var buf bytes.Buffer
		err = tmpl.Execute(&buf, params)
		if err != nil {
			panic(err)
		}
		return buf.String(), nil
	}
	tmpl, err := template.New("unimplementedConnectHandlerFile").Parse(unimplementedConnectHandlerFileTemplate)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, params)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

type connectHandlerFileParams struct {
	PackageName         string
	ServiceStructName   string
	ServicePackage      string
	ServicePackageAlias string
	RpcName             string
	RpcRequestName      string
	RpcResponseName     string
	RpcMethod           string
	MessageName         string
	MessageNameSingular string
}

func extractConnectHandlerFileParams(input ConnectHandlerFileInput, service *descriptorpb.ServiceDescriptorProto, method *descriptorpb.MethodDescriptorProto) connectHandlerFileParams {
	servicePackage := fmt.Sprintf("%s/%s", input.ApiPackage, strings.Replace(input.PackageName, ".", "/", -1))
	parts := helpers.SplitCamelCase(method.GetName())
	rpcMethod := parts[0]
	messageName := strings.TrimPrefix(method.GetName(), rpcMethod)
	messageName = strings.TrimSuffix(messageName, "Request")
	return connectHandlerFileParams{
		PackageName:         helpers.ToGoPackageName(input.PackageName),
		ServiceStructName:   helpers.ToLowerFirst(service.GetName()),
		ServicePackage:      servicePackage,
		ServicePackageAlias: helpers.ToGoAlias(input.PackageName),
		RpcName:             method.GetName(),
		RpcRequestName:      fmt.Sprintf("%sRequest", method.GetName()),
		RpcResponseName:     fmt.Sprintf("%sResponse", method.GetName()),
		RpcMethod:           rpcMethod,
		MessageName:         messageName,
		MessageNameSingular: strings.TrimSuffix(messageName, "s"),
	}
}

var getConnectHandlerFileTemplate = `
package {{.PackageName}}

import (
	{{.ServicePackageAlias}} "{{.ServicePackage}}"
	"connectrpc.com/connect"
	"context"
)

func (service {{.ServiceStructName}}) {{.RpcName}}(ctx context.Context, request *connect.Request[{{.ServicePackageAlias}}.{{.RpcRequestName}}]) (*connect.Response[{{.ServicePackageAlias}}.{{.RpcResponseName}}], error) {
	dbParams := {{.RpcRequestName}}ToDbParam(request.Msg)
	found, err := service.store.{{.RpcName}}(ctx, dbParams)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	response := connect.NewResponse(&{{.ServicePackageAlias}}.{{.RpcResponseName}}{
		{{.MessageName}}: DB{{.MessageNameSingular}}ToAPI{{.MessageNameSingular}}(found),
	})
	return response, nil
}
`

var listConnectHandlerFileTemplate = `
package {{.PackageName}}

import (
	"_shared/go-sdk/collections"
	{{.ServicePackageAlias}} "{{.ServicePackage}}"
	"connectrpc.com/connect"
	"context"
)

func (service {{.ServiceStructName}}) {{.RpcName}}(ctx context.Context, request *connect.Request[{{.ServicePackageAlias}}.{{.RpcRequestName}}]) (*connect.Response[{{.ServicePackageAlias}}.{{.RpcResponseName}}], error) {
	dbParams := {{.RpcRequestName}}ToDbParam(request.Msg)
	found, err := service.store.{{.RpcName}}(ctx, dbParams)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	response := connect.NewResponse(&{{.ServicePackageAlias}}.{{.RpcResponseName}}{
		{{.MessageName}}: collections.Map(found, DB{{.MessageNameSingular}}ToAPI{{.MessageNameSingular}}),
	})
	return response, nil
}
`

var createConnectHandlerFileTemplate = `
package {{.PackageName}}

import (
	{{.ServicePackageAlias}} "{{.ServicePackage}}"
	"connectrpc.com/connect"
	"context"
	"github.com/gofrs/uuid"
)

func (service {{.ServiceStructName}}) {{.RpcName}}(ctx context.Context, request *connect.Request[{{.ServicePackageAlias}}.{{.RpcRequestName}}]) (*connect.Response[{{.ServicePackageAlias}}.{{.RpcResponseName}}], error) {
	dbParams := {{.RpcRequestName}}ToDbParam(request.Msg)
	id := uuid.Must(uuid.NewV4())
	dbParams.ID = id
	created, err := service.store.{{.RpcName}}(ctx, dbParams)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	response := connect.NewResponse(&{{.ServicePackageAlias}}.{{.RpcResponseName}}{
		{{.MessageName}}: DB{{.MessageName}}ToAPI{{.MessageName}}(created),
	})
	return response, nil
}
`

var updateConnectHandlerFileTemplate = `
package {{.PackageName}}

import (
	{{.ServicePackageAlias}} "{{.ServicePackage}}"
	"connectrpc.com/connect"
	"context"
	"github.com/gofrs/uuid"
)

func (service {{.ServiceStructName}}) {{.RpcName}}(ctx context.Context, request *connect.Request[{{.ServicePackageAlias}}.{{.RpcRequestName}}]) (*connect.Response[{{.ServicePackageAlias}}.{{.RpcResponseName}}], error) {
	dbParams := {{.RpcRequestName}}ToDbParam(request.Msg)
	id := uuid.Must(uuid.NewV4())
	dbParams.ID = id
	updated, err := service.store.{{.RpcName}}(ctx, dbParams)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	response := connect.NewResponse(&{{.ServicePackageAlias}}.{{.RpcResponseName}}{
		{{.MessageName}}: DB{{.MessageName}}ToAPI{{.MessageName}}(updated),
	})
	return response, nil
}`

var unimplementedConnectHandlerFileTemplate = `
package {{.PackageName}}

import (
	"connectrpc.com/connect"
	"context"
	{{.ServicePackageAlias}} "{{.ServicePackage}}"
)

func (service {{.ServiceStructName}}) {{.RpcName}}(ctx context.Context, request *connect.Request[{{.ServicePackageAlias}}.{{.RpcRequestName}}]) (*connect.Response[{{.ServicePackageAlias}}.{{.RpcResponseName}}], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, nil)
}
`
