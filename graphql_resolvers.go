package vivographql

//https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1
import (
	"fmt"
	"log"

	//ge "github.com/OIT-ads-web/graphql_endpoint"
	//"github.com/OIT-ads-web/graphql_endpoint/elastic"
	"github.com/graphql-go/graphql"
	//"github.com/vivo-community/vivo-graphql/elastic"

	ms "github.com/mitchellh/mapstructure"
)

func personResolver(params graphql.ResolveParams) (interface{}, error) {
	id := params.Args["id"].(string)
	log.Printf("looking for person %s\n", id)

	person, err := FindPerson(id)
	return person, err
}

func publicationResolver(params graphql.ResolveParams) (interface{}, error) {
	id := params.Args["id"].(string)
	log.Printf("looking for publication %s\n", id)

	person, err := FindPublication(id)
	return person, err
}

func grantResolver(params graphql.ResolveParams) (interface{}, error) {
	id := params.Args["id"].(string)
	log.Printf("looking for grant %s\n", id)

	person, err := FindGrant(id)
	return person, err
}

// NOTE: this duplicates structure here:
// var PersonFilter *graphql.InputObject
// not sure best way to go about this
type CommonFilter struct {
	Limit  int
	Offset int
	Query  string
}

// NOTE: these aren't different now, but dealing with
// facets would probably make them different
type PersonFilterParam struct {
	Filter CommonFilter
}

type PublicationFilterParam struct {
	Filter CommonFilter
}

type GrantFilterParam struct {
	Filter CommonFilter
}

func convertPeopleFilter(params graphql.ResolveParams) (PersonFilterParam, error) {
	result := PersonFilterParam{
		Filter: CommonFilter{Limit: 100, Offset: 0, Query: ""},
	}
	err := ms.Decode(params.Args, &result)
	return result, err
}

func convertPublicationFilter(params graphql.ResolveParams) (PublicationFilterParam, error) {
	// default values?
	result := PublicationFilterParam{
		Filter: CommonFilter{Limit: 100, Offset: 0, Query: ""},
	}
	err := ms.Decode(params.Args, &result)
	return result, err
}

func convertGrantFilter(params graphql.ResolveParams) (GrantFilterParam, error) {
	// default values?
	result := GrantFilterParam{
		Filter: CommonFilter{Limit: 100, Offset: 0, Query: ""},
	}
	err := ms.Decode(params.Args, &result)
	return result, err
}

func peopleResolver(params graphql.ResolveParams) (interface{}, error) {
	// TODO: not finding a good way to default these
	// e.g. if filter is not sent at all
	limit := 100
	offset := 0
	query := ""
	filter, err := convertPeopleFilter(params)

	if err == nil {
		limit = filter.Filter.Limit
		offset = filter.Filter.Offset
		// NOTE: this is not that great
		query = fmt.Sprintf("*:%v*", filter.Filter.Query)
	}

	fmt.Printf("limit=%v,offset=%v,query=%v\n", limit, offset, query)
	personList, err := FindPeople(limit, offset, query)
	return personList, err
}

func publicationsResolver(params graphql.ResolveParams) (interface{}, error) {
	// TODO: not finding a good way to default these
	limit := 100
	offset := 0
	query := ""
	filter, err := convertPublicationFilter(params)
	if err == nil {
		limit = filter.Filter.Limit
		offset = filter.Filter.Offset
		query = fmt.Sprintf("*:%v*", filter.Filter.Query)
	}

	publications, err := FindPublications(limit, offset, query)
	return publications, err
}

func personPublicationResolver(params graphql.ResolveParams) (interface{}, error) {
	person, _ := params.Source.(Person)

	limit := params.Args["limit"].(int)
	offset := params.Args["offset"].(int)

	publicationList, err := FindPersonPublications(person.Id, limit, offset)
	return func() (interface{}, error) {
		return &publicationList, err
	}, nil
}

func grantsResolver(params graphql.ResolveParams) (interface{}, error) {
	limit := 100
	offset := 0
	query := ""
	filter, err := convertGrantFilter(params)
	if err == nil {
		limit = filter.Filter.Limit
		offset = filter.Filter.Offset
		query = fmt.Sprintf("*:%v*", filter.Filter.Query)
	}
	grants, err := FindGrants(limit, offset, query)
	return grants, err
}

func personGrantResolver(params graphql.ResolveParams) (interface{}, error) {
	person, _ := params.Source.(Person)

	limit := params.Args["limit"].(int)
	offset := params.Args["offset"].(int)

	grants, err := FindPersonGrants(person.Id, limit, offset)

	return func() (interface{}, error) {
		return &grants, err
	}, nil
}
