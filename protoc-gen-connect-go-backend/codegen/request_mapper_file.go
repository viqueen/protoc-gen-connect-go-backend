package codegen

import (
	"fmt"
	"github.com/viqueen/protoc-gen-connect-go-backend/protoc-gen-connect-go-backend/helpers"
	"google.golang.org/protobuf/types/descriptorpb"
	"strings"
	"text/template"
)

type RequestMapperFileInput struct {
	PackageName    string
	ApiPackage     string
	DataGenPackage string
}

func RequestMapperFile(input RequestMapperFileInput, message *descriptorpb.DescriptorProto) (string, error) {
	tmpl, err := template.New("requestMapperFile").Parse(requestMapperFileTemplate)
	if err != nil {
		return "", err
	}
	params := extractRequestMapperFileParams(input, message)
	var buf strings.Builder
	err = tmpl.Execute(&buf, params)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

type requestMapperFileParams struct {
	PackageName         string
	ServicePackageAlias string
	ServicePackage      string
	DataGenPackage      string
	RequestName         string
	RpcName             string
	Fields              []requestField
	HasIdField          bool
}

type requestField struct {
	ApiFieldName string
	DbFieldName  string
}

func extractRequestMapperFileParams(input RequestMapperFileInput, message *descriptorpb.DescriptorProto) requestMapperFileParams {
	servicePackage := fmt.Sprintf("%s/%s", input.ApiPackage, strings.Replace(input.PackageName, ".", "/", -1))
	var fields []requestField
	hasIdField := false
	for _, field := range message.GetField() {
		isID := strings.HasSuffix(field.GetName(), "_id") || field.GetName() == "id"
		if isID {
			hasIdField = true
		}
		goFieldName := helpers.ToGoFieldName(field.GetName())
		camelCaseFieldName := helpers.SnakeToCamel(field.GetName())
		apiFieldName := fmt.Sprintf("request.Get%s()", camelCaseFieldName)
		if isID {
			apiFieldName = fmt.Sprintf("uuid.FromStringOrNil(request.Get%s())", camelCaseFieldName)
		}
		fields = append(fields, requestField{
			ApiFieldName: apiFieldName,
			DbFieldName:  goFieldName,
		})
	}
	return requestMapperFileParams{
		PackageName:         helpers.ToGoPackageName(input.PackageName),
		ServicePackageAlias: helpers.ToGoAlias(input.PackageName),
		ServicePackage:      servicePackage,
		DataGenPackage:      input.DataGenPackage,
		RequestName:         message.GetName(),
		RpcName:             strings.TrimSuffix(message.GetName(), "Request"),
		Fields:              fields,
		HasIdField:          hasIdField,
	}
}

var requestMapperFileTemplate = `
package {{ .PackageName }}

import (
	{{.ServicePackageAlias}} "{{.ServicePackage}}"
	{{ if .HasIdField }}"github.com/gofrs/uuid"{{end}}
    gendb "{{.DataGenPackage}}"
)

func {{ .RequestName }}ToDbParam(request *{{.ServicePackageAlias}}.{{.RequestName}}) *gendb.{{.RpcName}}Params {
	if request == nil {
		return nil
	}
	return &gendb.{{.RpcName}}Params{
		{{range .Fields}}{{.DbFieldName}}: {{.ApiFieldName}},
		{{end}}
	}
}
`
