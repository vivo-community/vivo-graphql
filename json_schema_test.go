package vivographql_test

import (
	"fmt"
	"testing"

	vq "github.com/vivo-community/vivo-graphql"
)

func TestInValidPerson(t *testing.T) {
	// best way to get this value in?
	//schemaDir := os.Getenv("SCHEMAS_DIR")
	schemaDir := "json-schemas/"
	var conf = vq.Config{Schemas: vq.SchemasPath{Dir: schemaDir}}
	vq.LoadJsonSchemas(conf)

	json := `{
		"id":"per000001"
	}`
	schema := vq.RetrieveJsonSchema("person")
	valid, _ := vq.Validate(schema, json)
	// TODO: should probably check the error list
	if valid {
		t.Error(fmt.Printf("should be NOT valid %v\n", json))
	}
}

func TestValidPerson(t *testing.T) {
	// FIXME: best way to get this value in?
	schemaDir := "json-schemas/"
	var conf = vq.Config{Schemas: vq.SchemasPath{Dir: schemaDir}}
	vq.LoadJsonSchemas(conf)

	json := `{
		"id":"per000001",
		"name": {
			"firstName": "Lester",
			"lastName": "Tester"
		}
	}`
	schema := vq.RetrieveJsonSchema("person")
	valid, errors := vq.Validate(schema, json)
	for error := range errors {
		fmt.Printf("%v\n", error)
	}
	// TODO: should check the error list
	if !valid {
		t.Error(fmt.Printf("should be valid %v\n", json))
	}
}
