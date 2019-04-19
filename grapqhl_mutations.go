package vivographql

import (
	"errors"

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
	// do something like this:
	// indexer.SavePerson(newPerson)

	// this might make it async
	return func() (interface{}, error) {
		return &newPerson, err
	}, nil
}

// TODO: these both will have to validate - so should share that logic
func personValidation(params graphql.ResolveParams) (interface{}, error) {
	// marshall and cast the argument value
	id, _ := params.Args["id"].(string)

	// perform mutation operation here
	// for e.g. create a Person and save to Indexer.
	newPerson := &Person{
		Id: id,
	}
	err := errors.New("not implemented")
	// do something like this:
	// indexer.SavePerson(newPerson)

	// this might make it async
	return func() (interface{}, error) {
		return &newPerson, err
	}, nil
}
