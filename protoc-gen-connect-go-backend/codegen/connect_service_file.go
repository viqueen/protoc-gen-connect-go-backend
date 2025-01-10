package codegen

import (
	"bytes"
	"fmt"
	"github.com/viqueen/protoc-gen-connect-go-backend/protoc-gen-connect-go-backend/helpers"
	"google.golang.org/protobuf/types/descriptorpb"
	"strings"
	"text/template"
)

type ConnectServiceFileInput struct {
	PackageName    string
	ApiPackage     string
	DataGenPackage string
}

func ConnectServiceFile(input ConnectServiceFileInput, service *descriptorpb.ServiceDescriptorProto) (string, error) {
	tmpl, err := template.New("connectServiceFile").Parse(connectServiceFileTemplate)
	if err != nil {
		return "", err
	}
	params := extractConnectServiceFileParams(input, service)
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, params)
	if err != nil {
		panic(err)
	}
	return buf.String(), nil
}

type connectServiceFileParams struct {
	PackageName                string
	ServiceName                string
	ServiceStructName          string
	ServicePackage             string
	ServiceConnectPackageAlias string
	ServiceConnectPackage      string
	DataGenPackage             string
}

func extractConnectServiceFileParams(input ConnectServiceFileInput, service *descriptorpb.ServiceDescriptorProto) connectServiceFileParams {
	servicePackage := fmt.Sprintf("%s/%s", input.ApiPackage, strings.Replace(input.PackageName, ".", "/", -1))
	serviceConnectPackageAlias := fmt.Sprintf("%sconnect", helpers.ToGoAlias(input.PackageName))
	serviceConnectPackage := fmt.Sprintf("%s/%s", servicePackage, serviceConnectPackageAlias)
	return connectServiceFileParams{
		PackageName:                helpers.ToGoPackageName(input.PackageName),
		ServiceName:                service.GetName(),
		ServiceStructName:          helpers.ToLowerFirst(service.GetName()),
		ServicePackage:             servicePackage,
		ServiceConnectPackageAlias: serviceConnectPackageAlias,
		ServiceConnectPackage:      serviceConnectPackage,
		DataGenPackage:             input.DataGenPackage,
	}
}

var connectServiceFileTemplate = `
package {{.PackageName}}

import (
	"{{.ServiceConnectPackage}}"
	gendb "{{.DataGenPackage}}"
)

type {{.ServiceStructName}} struct {
	store gendb.Querier
}

type {{.ServiceName}}Config struct {
	Store gendb.Querier
}

func New{{.ServiceName}}(config {{.ServiceName}}Config) {{.ServiceConnectPackageAlias}}.{{.ServiceName}}Handler {
	return &{{.ServiceStructName}}{
		store: config.Store,
	}
}
`
