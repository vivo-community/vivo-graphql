package vivographql

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/graphql-go/graphql"
)

// TODO: much more to do here
func personMutation(params graphql.ResolveParams) (interface{}, error) {
	// marshall and cast the argument value
	id, _ := params.Args["id"].(string)

	// perform mutation operation here
	// for e.g. create a Person and save to Indexer.
	newPerson := &Person{
		Id: id,
	}
	err := errors.New("not implemented")
	// would do something like this:
	// indexer.SavePerson(newPerson)

	// this might make it async
	return func() (interface{}, error) {
		return &newPerson, err
	}, nil
}

// TODO: these both will have to validate - so should share that logic
func personValidation(params graphql.ResolveParams) (interface{}, error) {
	id, _ := params.Args["id"].(string)

	// perform mutation operation here
	// for e.g. create a Person and save to Indexer.
	newPerson := &Person{
		Id: id,
	}

	schema := RetrieveJsonSchema("person")
	b, err := json.Marshal(newPerson)

	valid, errList := Validate(schema, string(b))

	if !valid {
		// ResultError object
		for _, e := range errList {
			// TODO: should do something better than this
			fmt.Printf("%s\n", e.Description())
		}
		// way to return multiple errors? errList ???
		err = errors.New("not valid")
	}

	//err = errors.New("not implemented")
	// this might make it async
	return func() (interface{}, error) {
		return &newPerson, err
	}, nil
}
