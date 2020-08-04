package main

import (
	"fmt"
	"go-sso/cli"
	"go-sso/db/inter"
	"go-sso/di"
	"log"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
// @BasePath /v2
func main() {
	container := di.BuildContainer()
	err := container.Invoke(di.PrintConfig)
	if err != nil {
		log.Fatal(err)
	}
	cli.Execute()
}
