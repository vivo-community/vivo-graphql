package vivographql

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

/* usage:

schema := vivographql.RetrieveSchema("person")
json := `{some_data: "hello"}`
valid := vivographql.Validate(schema, string(json))

*/

var schemas map[string]*gojsonschema.Schema

// preloading at start - then storing by key typeName
func LoadSchemas(conf Config) {
	schemas = make(map[string]*gojsonschema.Schema)

	log.Printf("looking for schemas in %s\n", conf.Schemas.Dir)
	schemaFiles, err := filepath.Glob(conf.Schemas.Dir + "*.json")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range schemaFiles {
		log.Printf("trying to parse schema %s\n", file)
		b, err := ioutil.ReadFile(file) // just pass the file name
		if err != nil {
			fmt.Print(err)
		}

		schemaDef := string(b)
		loader1 := gojsonschema.NewStringLoader(schemaDef)
		schema, err := gojsonschema.NewSchema(loader1)

		if err != nil {
			fmt.Println("could not load schema")
			panic(err)
		}

		fileName := filepath.Base(file)

		typeName := strings.Replace(fileName, ".schema.json", "", 1)
		// store as key typeName
		log.Printf("putting schema in cache[%s]\n", typeName)
		schemas[typeName] = schema
	}

}

func RetrieveSchema(typeName string) *gojsonschema.Schema {
	switch typeName {
	case "person":
		return schemas["person"]
	case "publication":
		return schemas["publication"]
	case "grant":
		return schemas["grant"]
	case "funding-role":
		return schemas["funding-role"]
	case "authorship":
		return schemas["authorship"]
	case "affiliation":
		return schemas["affiliation"]
	case "education":
		return schemas["education"]
	default:
		err := fmt.Sprintf("could not load schema, cancelling %s\n", typeName)
		panic(err)
	}
}

// maybe allow []byte for data param
func Validate(schema *gojsonschema.Schema, data string) bool {
	docLoader := gojsonschema.NewStringLoader(data)
	result, err := schema.Validate(docLoader)

	if err != nil {
		fmt.Println("error validating")
		return false
	}

	if result.Valid() {
		fmt.Printf("The document is valid\n")
		if err != nil {
			fmt.Printf("- %s\n", err)
		}
		return true
	} else {
		fmt.Printf("The document is not valid. see errors :\n")
		for _, err := range result.Errors() {
			// Err implements the ResultError interface
			fmt.Printf("- %s\n", err)
		}
		return false
	}
}
