package vivographql

import "github.com/graphql-go/graphql"

// TODO: much more to do here
func personMutation(params graphql.ResolveParams) (interface{}, error) {
	// marshall and cast the argument value
	id, _ := params.Args["id"].(string)

	// perform mutation operation here
	// for e.g. create a Person and save to Indexer.
	newPerson := &Person{
		Id: id,
	}
	// indexer.SavePerson(newPerson)
	return newPerson, nil
}
