package main

import "github.com/philiplambok/task-api/cmd"

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=config/oapi-gen.yml api/openapi.yml
func main() {
	cmd.Execute()
}
