package main

import (
	"log"
	"os"

	"vegeta-server/restapi"
	"vegeta-server/restapi/operations"

	loads "github.com/go-openapi/loads"
	goflags "github.com/jessevdk/go-flags"
)

func main() {

	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewVegetaAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer server.Shutdown() // nolint: errcheck

	parser := goflags.NewParser(server, goflags.Default)
	parser.ShortDescription = "Vegeta REST API"
	parser.LongDescription = "This is a RESTful API for the vegeta load-testing utility. Vegeta is a versatile HTTP load testing tool built out of a need to drill HTTP services with a constant request rate.\n"

	server.ConfigureFlags()
	for _, optsGroup := range api.CommandLineOptionsGroups {
		_, err := parser.AddGroup(optsGroup.ShortDescription, optsGroup.LongDescription, optsGroup.Options)
		if err != nil {
			log.Fatalln(err)
		}
	}

	if _, err := parser.Parse(); err != nil {
		code := 1
		if fe, ok := err.(*goflags.Error); ok {
			if fe.Type == goflags.ErrHelp {
				code = 0
			}
		}
		os.Exit(code)
	}

	server.ConfigureAPI()

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
