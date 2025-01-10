package codegen

import (
	"bytes"
	"fmt"
	"github.com/viqueen/protoc-gen-connect-go-backend/protoc-gen-connect-go-backend/helpers"
	"google.golang.org/protobuf/types/descriptorpb"
	"strings"
	"text/template"
)

type ConnectServerFileInput struct {
	DataGenPackage  string
	PackageName     string
	InternalPackage string
	ApiPackage      string
}

func ConnectServerFile(input ConnectServerFileInput, services []*descriptorpb.ServiceDescriptorProto) (string, error) {
	tmpl, err := template.New("connectServerFile").Parse(connectServerFileTemplate)
	if err != nil {
		return "", err
	}
	params := extractConnectServerFileParams(input, services)
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, params)
	if err != nil {
		panic(err)
	}
	return buf.String(), nil
}

type connectServerFileParams struct {
	Services       []serviceDef
	DataGenPackage string
}

type serviceDef struct {
	ServiceNameLower            string
	ServiceInternalPackage      string
	ServiceInternalPackageAlias string
	ServiceName                 string
	ServiceConnectPackage       string
	ServiceConnectPackageAlias  string
}

func extractConnectServerFileParams(input ConnectServerFileInput, services []*descriptorpb.ServiceDescriptorProto) connectServerFileParams {
	var serviceDefs []serviceDef
	for _, service := range services {
		servicePackage := fmt.Sprintf("%s/%s", input.ApiPackage, strings.Replace(input.PackageName, ".", "/", -1))
		serviceConnectPackageAlias := fmt.Sprintf("%sconnect", helpers.ToGoAlias(input.PackageName))
		serviceConnectPackage := fmt.Sprintf("%s/%s", servicePackage, serviceConnectPackageAlias)
		serviceDefs = append(serviceDefs, serviceDef{
			ServiceNameLower:            strings.TrimSuffix(strings.ToLower(service.GetName()), "service"),
			ServiceName:                 strings.TrimSuffix(service.GetName(), "Service"),
			ServiceInternalPackage:      fmt.Sprintf("%s/api-%s", input.InternalPackage, strings.ReplaceAll(input.PackageName, ".", "-")),
			ServiceInternalPackageAlias: strings.ReplaceAll(input.PackageName, ".", ""),
			ServiceConnectPackage:       serviceConnectPackage,
			ServiceConnectPackageAlias:  serviceConnectPackageAlias,
		})
	}
	return connectServerFileParams{
		Services:       serviceDefs,
		DataGenPackage: input.DataGenPackage,
	}
}

var connectServerFileTemplate = `
package main

import (
	connectcors "connectrpc.com/cors"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"log"
	"net/http"
	gendb "{{.DataGenPackage}}"
	{{range .Services}}
	{{.ServiceInternalPackageAlias}} "{{.ServiceInternalPackage}}"
	{{.ServiceConnectPackageAlias}} "{{.ServiceConnectPackage}}"
	{{end}}

	"connectrpc.com/otelconnect"
)

func main() {
	otelInterceptor, otelErr := otelconnect.NewInterceptor()
	if otelErr != nil {
		log.Fatalf("failed to initialise otel interceptor: %v", otelErr)
	}

	db, dbErr := initialiseDB()
	if dbErr != nil {
		log.Fatalf("failed to initialise db: %v", dbErr)
	}
	dataStore := gendb.New(db)

	mux := http.NewServeMux()
	{{range .Services}}
	{{.ServiceNameLower}}Service := {{.ServiceInternalPackageAlias}}.New{{.ServiceName}}Service({{.ServiceInternalPackageAlias}}.{{.ServiceName}}ServiceConfig{Store: dataStore})
	{{.ServiceNameLower}}Path, {{.ServiceNameLower}}Handler := {{.ServiceConnectPackageAlias}}.New{{.ServiceName}}ServiceHandler({{.ServiceNameLower}}Service, connect.WithInterceptors(otelInterceptor))
	mux.Handle({{.ServiceNameLower}}Path, {{.ServiceNameLower}}Handler)
	{{end}}

	h2cMux := h2c.NewHandler(mux, &http2.Server{})
	log.Printf("starting server on :8080")
	err := http.ListenAndServe(":8080", withCORS(h2cMux))
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func initialiseDB() (*sql.DB, error) {
	connectionString := "postgres://canopy:canopy@localhost:5432/canopy?sslmode=disable"
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func withCORS(h http.Handler) http.Handler {
	middleware := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: connectcors.AllowedMethods(),
		AllowedHeaders: connectcors.AllowedHeaders(),
		ExposedHeaders: connectcors.ExposedHeaders(),
	})
	return middleware.Handler(h)
}
`
