package codegen

import (
	"fmt"
	"google.golang.org/protobuf/types/descriptorpb"
	"strings"
	"text/template"
)

type DataMapperFileInput struct {
	PackageName    string
	ApiPackage     string
	DataGenPackage string
}

func DataMapperFile(input DataMapperFileInput, message *descriptorpb.DescriptorProto) (string, error) {
	tmpl, err := template.New("dataMapperFile").Parse(dataMapperFileTemplate)
	if err != nil {
		return "", err
	}
	params := extractDataMapperFileParams(input, message)
	var buf strings.Builder
	err = tmpl.Execute(&buf, params)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

type dataMapperFileParams struct {
	PackageName         string
	ServicePackageAlias string
	ServicePackage      string
	DataGenPackage      string
	DbTypeName          string
	MessageName         string
	Fields              []dataField
}

type dataField struct {
	ApiFieldName string
	DbFieldName  string
}

func extractDataMapperFileParams(input DataMapperFileInput, message *descriptorpb.DescriptorProto) dataMapperFileParams {
	servicePackage := fmt.Sprintf("%s/%s", input.ApiPackage, strings.Replace(input.PackageName, ".", "/", -1))
	var fields []dataField
	for _, field := range message.GetField() {
		isID := strings.HasSuffix(field.GetName(), "_id") || field.GetName() == "id"
		goFieldName := toGoFieldName(field.GetName())
		dbFieldName := fmt.Sprintf("input.%s", goFieldName)
		if isID {
			dbFieldName = fmt.Sprintf("input.%s.String()", goFieldName)
		}
		fields = append(fields, dataField{
			ApiFieldName: strings.Title(field.GetName()),
			DbFieldName:  dbFieldName,
		})
	}
	return dataMapperFileParams{
		PackageName:         toGoPackageName(input.PackageName),
		ServicePackageAlias: toGoAlias(input.PackageName),
		ServicePackage:      servicePackage,
		DataGenPackage:      input.DataGenPackage,
		DbTypeName:          message.GetName(),
		MessageName:         message.GetName(),
		Fields:              fields,
	}
}

var dataMapperFileTemplate = `
package {{.PackageName}}

import (
	{{.ServicePackageAlias}} "{{.ServicePackage}}"
	gendb "{{.DataGenPackage}}"
)

func DB{{.MessageName}}ToAPI{{.MessageName}}(input *gendb.{{.DbTypeName}}) *{{.ServicePackageAlias}}.{{.MessageName}} {
	if input == nil {
		return nil
	}
	return &{{.ServicePackageAlias}}.{{.MessageName}}{
		{{range .Fields}}{{.ApiFieldName}}: {{.DbFieldName}},
		{{end}}
	}
}
`
