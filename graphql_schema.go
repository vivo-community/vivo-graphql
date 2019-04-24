package vivographql

import (
	"sync"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

func MakeGraphqlHandler() *handler.Handler {
	schema := GetVivoSchema()
	h := handler.New(&handler.Config{
		Schema:   schema,
		GraphiQL: true,
		Pretty:   true,
	})
	return h
}

// stole 'once' idea from here:
// http://marcio.io/2015/07/singleton-pattern-in-go/
var schema *graphql.Schema

// NOTE: this is a little weird - only
// one 'once' per package
var oneSchema sync.Once

func GetVivoSchema() *graphql.Schema {
	oneSchema.Do(func() {
		schema = MakeSchema()
	})
	return schema
}

func MakeSchema() *graphql.Schema {
	// ignoring error
	var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query:    RootQuery,
		Mutation: RootMutation,
	})
	return &schema
}

// how to do something like this?
//https://medium.com/@bradford_hamilton/building-an-api-with-graphql-and-go-9350df5c9356
/*
type Root struct {
	Query *graphql.Object
}

type Resolver struct {
	idx *Indexer
}

func NewRoot(indexer *Indexer) *Root {
	resolver := Resolver{idx: indexer}
	root := Root{Query: RootQuery}
	return &root
}

func (r *Resolver) PersonResolver(p graphql.ResolveParams) (interface{}, error) {
   // etc...
}

*/

var RootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"personList":      GetPeople,
		"person":          GetPerson,
		"grant":           GetGrant,
		"publication":     GetPublication,
		"publicationList": GetPublications,
		"grantList":       GetGrants,
	},
})

var RootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{
		"createPerson": CreatePerson,
		//"createGrant": CreateGrant,
		//"createPublication": CreatePublication,
		// https://github.com/graphql-go/graphql/issues/234
		//"addAffilition": AddAffilition,
		//"addEducation": AddEducation,
		//"createCourse": CreateCourse,
		// etc...
		"validatePerson": ValidatePerson,
		//"validatePublication": ValidatePublication,
		// etc...
	},
})

// TODO: probably need a complex input type e.g.
var PersonInput *graphql.InputObject = graphql.NewInputObject(graphql.InputObjectConfig{
	Name:        "PersonInput",
	Description: "Person to Add/Update",
	Fields: graphql.InputObjectConfigFieldMap{
		"id": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
	},
})

var CreatePerson = &graphql.Field{
	Type: person, // the return type for this field
	Args: graphql.FieldConfigArgument{
		// TODO: all fields here?
		"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
	},
	// TODO: maybe something like this? not sure
	//Args: graphql.FieldConfigArgument{
	//	"person": &graphql.ArgumentConfig{Type: PersonInput},
	//},
	Resolve: personMutation,
}

var ValidatePerson = &graphql.Field{
	Type: person, // the return type for this field
	Args: graphql.FieldConfigArgument{
		// TODO: all fields here?
		"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
	},
	// TODO: maybe something like this? not sure
	//Args: graphql.FieldConfigArgument{
	//	"person": &graphql.ArgumentConfig{Type: PersonInput},
	//},
	Resolve: personValidation,
}

var GetPerson = &graphql.Field{
	Type:        person,
	Description: "Get Person",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
	},
	// how to switch resolver based on config?  solr/elastic ?
	// also how to get Resolver object into every Field?
	//Resolve: resolver.PersonResolver
	Resolve: personResolver,
}

var GetPublication = &graphql.Field{
	Type:        publication,
	Description: "Get Publication",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
	},
	Resolve: publicationResolver,
}

var GetGrant = &graphql.Field{
	Type:        grant,
	Description: "Get Grant",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
	},
	Resolve: grantResolver,
}

var PersonFilter *graphql.InputObject = graphql.NewInputObject(graphql.InputObjectConfig{
	Name:        "PersonFilter",
	Description: "Filter on People List",
	Fields: graphql.InputObjectConfigFieldMap{
		"limit": &graphql.InputObjectFieldConfig{
			Type:         graphql.Int,
			DefaultValue: 100,
		},
		"offset": &graphql.InputObjectFieldConfig{
			Type:         graphql.Int,
			DefaultValue: 0,
		},
		"query": &graphql.InputObjectFieldConfig{
			Type:         graphql.String,
			DefaultValue: "",
		},
	},
})

var GetPeople = &graphql.Field{
	Type:        personList,
	Description: "Get all people",
	Args: graphql.FieldConfigArgument{
		"filter": &graphql.ArgumentConfig{Type: PersonFilter},
	},
	Resolve: peopleResolver,
}

// TODO: very likely a way to avoid the code duplication
var PublicationFilter *graphql.InputObject = graphql.NewInputObject(graphql.InputObjectConfig{
	Name:        "PublicationFilter",
	Description: "Filter on Publication List",
	Fields: graphql.InputObjectConfigFieldMap{
		"limit": &graphql.InputObjectFieldConfig{
			Type:         graphql.Int,
			DefaultValue: 100,
		},
		"offset": &graphql.InputObjectFieldConfig{
			Type:         graphql.Int,
			DefaultValue: 0,
		},
		"query": &graphql.InputObjectFieldConfig{
			Type:         graphql.String,
			DefaultValue: "",
		},
	},
})

var GetPublications = &graphql.Field{
	Type:        publicationList,
	Description: "Get all publications",
	Args: graphql.FieldConfigArgument{
		"filter": &graphql.ArgumentConfig{Type: PublicationFilter},
	},
	Resolve: publicationsResolver,
}

var GrantFilter *graphql.InputObject = graphql.NewInputObject(graphql.InputObjectConfig{
	Name:        "GrantFilter",
	Description: "Filter on Grant List",
	Fields: graphql.InputObjectConfigFieldMap{
		"limit": &graphql.InputObjectFieldConfig{
			Type:         graphql.Int,
			DefaultValue: 100,
		},
		"offset": &graphql.InputObjectFieldConfig{
			Type:         graphql.Int,
			DefaultValue: 0,
		},
		"query": &graphql.InputObjectFieldConfig{
			Type:         graphql.String,
			DefaultValue: "",
		},
	},
})

var GetGrants = &graphql.Field{
	Type:        grantList,
	Description: "Get all grants",
	Args: graphql.FieldConfigArgument{
		"filter": &graphql.ArgumentConfig{Type: GrantFilter},
	},
	Resolve: grantsResolver,
}

// NOTE: removed Type from end of names per
// https://dave.cheney.net/practical-go/presentations/qcon-china.html
var pageInfo = graphql.NewObject(graphql.ObjectConfig{
	Name: "PageInfo",
	Fields: graphql.Fields{
		"perPage":    &graphql.Field{Type: graphql.Int},
		"page":       &graphql.Field{Type: graphql.Int},
		"totalPages": &graphql.Field{Type: graphql.Int},
		"count":      &graphql.Field{Type: graphql.Int},
	},
})

var grant = graphql.NewObject(graphql.ObjectConfig{
	Name: "Grant",
	Fields: graphql.Fields{
		"id":        &graphql.Field{Type: graphql.String},
		"label":     &graphql.Field{Type: graphql.String},
		"roleName":  &graphql.Field{Type: graphql.String},
		"startDate": &graphql.Field{Type: dateResolution},
		"endDate":   &graphql.Field{Type: dateResolution},
	},
})

var organization = graphql.NewObject(graphql.ObjectConfig{
	Name: "Organization",
	Fields: graphql.Fields{
		"id":    &graphql.Field{Type: graphql.String},
		"label": &graphql.Field{Type: graphql.String},
	},
})

var education = graphql.NewObject(graphql.ObjectConfig{
	Name: "Education",
	Fields: graphql.Fields{
		"credential":             &graphql.Field{Type: graphql.String},
		"credentialAbbreviation": &graphql.Field{Type: graphql.String},
		"organization":           &graphql.Field{Type: organization},
		"dateReceived":           &graphql.Field{Type: dateResolution},
	},
})

var affiliation = graphql.NewObject(graphql.ObjectConfig{
	Name: "Affiliation",
	Fields: graphql.Fields{
		"id":        &graphql.Field{Type: graphql.String},
		"label":     &graphql.Field{Type: graphql.String},
		"startDate": &graphql.Field{Type: dateResolution},
	},
})

var keyword = graphql.NewObject(graphql.ObjectConfig{
	Name: "Keyword",
	Fields: graphql.Fields{
		"uri":   &graphql.Field{Type: graphql.String},
		"label": &graphql.Field{Type: graphql.String},
	},
})

var extension = graphql.NewObject(graphql.ObjectConfig{
	Name: "Extension",
	Fields: graphql.Fields{
		"key":   &graphql.Field{Type: graphql.String},
		"value": &graphql.Field{Type: graphql.String},
	},
})

var dateResolution = graphql.NewObject(graphql.ObjectConfig{
	Name: "DateResolution",
	Fields: graphql.Fields{
		"dateTime":   &graphql.Field{Type: graphql.String},
		"resolution": &graphql.Field{Type: graphql.String},
	},
})

var personName = graphql.NewObject(graphql.ObjectConfig{
	Name: "PersonName",
	Fields: graphql.Fields{
		"firstName":  &graphql.Field{Type: graphql.String},
		"lastName":   &graphql.Field{Type: graphql.String},
		"middleName": &graphql.Field{Type: graphql.String},
	},
})

var personImage = graphql.NewObject(graphql.ObjectConfig{
	Name: "PersonImage",
	Fields: graphql.Fields{
		"main":      &graphql.Field{Type: graphql.String},
		"thumbnail": &graphql.Field{Type: graphql.String},
	},
})

var typeOf = graphql.NewObject(graphql.ObjectConfig{
	Name: "Type",
	Fields: graphql.Fields{
		"code":  &graphql.Field{Type: graphql.String},
		"label": &graphql.Field{Type: graphql.String},
	},
})

var personIdentifier = graphql.NewObject(graphql.ObjectConfig{
	Name: "PersonIdentifier",
	Fields: graphql.Fields{
		"orchid": &graphql.Field{Type: graphql.String},
		"isni":   &graphql.Field{Type: graphql.String},
	},
})

var personContact = graphql.NewObject(graphql.ObjectConfig{
	Name: "Contact",
	Fields: graphql.Fields{
		"emailList":    &graphql.Field{Type: graphql.NewList(email)},
		"locationList": &graphql.Field{Type: graphql.NewList(location)},
		"phoneList":    &graphql.Field{Type: graphql.NewList(phone)},
		"websiteList":  &graphql.Field{Type: graphql.NewList(website)},
	},
})

var serviceRole = graphql.NewObject(graphql.ObjectConfig{
	Name: "ServiceRole",
	Fields: graphql.Fields{
		"uri":          &graphql.Field{Type: graphql.String},
		"label":        &graphql.Field{Type: graphql.String},
		"description":  &graphql.Field{Type: graphql.String},
		"startDate":    &graphql.Field{Type: dateResolution},
		"endDate":      &graphql.Field{Type: dateResolution},
		"organization": &graphql.Field{Type: organization},
		"type":         &graphql.Field{Type: typeOf},
	},
})

var email = graphql.NewObject(graphql.ObjectConfig{
	Name: "Email",
	Fields: graphql.Fields{
		"uri":   &graphql.Field{Type: graphql.String},
		"label": &graphql.Field{Type: graphql.String},
		"type":  &graphql.Field{Type: typeOf},
	},
})

var phone = graphql.NewObject(graphql.ObjectConfig{
	Name: "Phone",
	Fields: graphql.Fields{
		"uri":   &graphql.Field{Type: graphql.String},
		"label": &graphql.Field{Type: graphql.String},
		"type":  &graphql.Field{Type: typeOf},
	},
})

var location = graphql.NewObject(graphql.ObjectConfig{
	Name: "Location",
	Fields: graphql.Fields{
		"uri":   &graphql.Field{Type: graphql.String},
		"label": &graphql.Field{Type: graphql.String},
		"type":  &graphql.Field{Type: typeOf},
	},
})

var website = graphql.NewObject(graphql.ObjectConfig{
	Name: "Website",
	Fields: graphql.Fields{
		"uri":   &graphql.Field{Type: graphql.String},
		"url":   &graphql.Field{Type: graphql.String},
		"label": &graphql.Field{Type: graphql.String},
		"type":  &graphql.Field{Type: typeOf},
	},
})

var courseTaught = graphql.NewObject(graphql.ObjectConfig{
	Name: "CourseTaught",
	Fields: graphql.Fields{
		"uri":          &graphql.Field{Type: graphql.String},
		"subject":      &graphql.Field{Type: graphql.String},
		"role":         &graphql.Field{Type: graphql.String},
		"courseName":   &graphql.Field{Type: graphql.String},
		"courseNumber": &graphql.Field{Type: graphql.String},
		"startDate":    &graphql.Field{Type: dateResolution},
		"endDate":      &graphql.Field{Type: dateResolution},
		"organization": &graphql.Field{Type: organization},
		"type":         &graphql.Field{Type: typeOf},
	},
})

var publicationIdentifier = graphql.NewObject(graphql.ObjectConfig{
	Name: "PublicationIdentifier",
	Fields: graphql.Fields{
		"isbn10": &graphql.Field{Type: graphql.String},
		"isbn13": &graphql.Field{Type: graphql.String},
		"pmid":   &graphql.Field{Type: graphql.String},
		"doi":    &graphql.Field{Type: graphql.String},
		"pmcid":  &graphql.Field{Type: graphql.String},
	},
})

var publicationVenue = graphql.NewObject(graphql.ObjectConfig{
	Name: "PublicationVenue",
	Fields: graphql.Fields{
		"uri":   &graphql.Field{Type: graphql.String},
		"label": &graphql.Field{Type: graphql.String},
	},
})

var publication = graphql.NewObject(graphql.ObjectConfig{
	Name: "Publication",
	Fields: graphql.Fields{
		"id":               &graphql.Field{Type: graphql.String},
		"uri":              &graphql.Field{Type: graphql.String},
		"title":            &graphql.Field{Type: graphql.String},
		"identifier":       &graphql.Field{Type: publicationIdentifier},
		"type":             &graphql.Field{Type: typeOf},
		"authorList":       &graphql.Field{Type: graphql.String},
		"venue":            &graphql.Field{Type: publicationVenue},
		"abstract":         &graphql.Field{Type: graphql.String},
		"dateStandardized": &graphql.Field{Type: dateResolution},
		"dateDisplay":      &graphql.Field{Type: graphql.String},
		"pageRange":        &graphql.Field{Type: graphql.String},
		"pageStart":        &graphql.Field{Type: graphql.String},
		"pageEnd":          &graphql.Field{Type: graphql.String},
		"issue":            &graphql.Field{Type: graphql.String},
		"volume":           &graphql.Field{Type: graphql.String},
	},
})

var overview = graphql.NewObject(graphql.ObjectConfig{
	Name: "Overview",
	Fields: graphql.Fields{
		"code":  &graphql.Field{Type: graphql.String},
		"label": &graphql.Field{Type: graphql.String},
	},
})

var authorship = graphql.NewObject(graphql.ObjectConfig{
	Name: "Authorship",
	Fields: graphql.Fields{
		"id":            &graphql.Field{Type: graphql.String},
		"uri":           &graphql.Field{Type: graphql.String},
		"publicationId": &graphql.Field{Type: graphql.String},
		"personId":      &graphql.Field{Type: graphql.String},
		"label":         &graphql.Field{Type: graphql.String},
	},
})

var person = graphql.NewObject(graphql.ObjectConfig{
	Name: "Person",
	Fields: graphql.Fields{
		"uri":              &graphql.Field{Type: graphql.String},
		"id":               &graphql.Field{Type: graphql.String},
		"sourceId":         &graphql.Field{Type: graphql.String},
		"primaryTitle":     &graphql.Field{Type: graphql.String},
		"name":             &graphql.Field{Type: personName},
		"image":            &graphql.Field{Type: personImage},
		"identifier":       &graphql.Field{Type: personIdentifier},
		"contact":          &graphql.Field{Type: personContact},
		"type":             &graphql.Field{Type: typeOf},
		"overviewList":     &graphql.Field{Type: graphql.NewList(overview)},
		"keywordList":      &graphql.Field{Type: graphql.NewList(keyword)},
		"extensionList":    &graphql.Field{Type: graphql.NewList(extension)},
		"affliationList":   &graphql.Field{Type: graphql.NewList(affiliation)},
		"educationList":    &graphql.Field{Type: graphql.NewList(education)},
		"serviceRoleList":  &graphql.Field{Type: graphql.NewList(serviceRole)},
		"courseTaughtList": &graphql.Field{Type: graphql.NewList(courseTaught)},
		// these can be paged, since they involve further queries
		"publicationList": &graphql.Field{
			Type: publicationList,
			// TODO: would probably want these to be like
			// PublicationList (e.g. a Filter object)
			Args: graphql.FieldConfigArgument{
				"limit":  &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 100},
				"offset": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 1},
			},
			Resolve: personPublicationResolver,
		},
		"grantList": &graphql.Field{
			Type: grantList,
			Args: graphql.FieldConfigArgument{
				"limit":  &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 100},
				"offset": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 1},
			},
			Resolve: personGrantResolver,
		},
	},
})

/*
FIXME: not sure best way to send in sorting parameters right now

also how to add enum of choices (e.g. asc, desc)
*/
var sortField = graphql.NewObject(graphql.ObjectConfig{
	Name: "SortField",
	Fields: graphql.Fields{
		"name":  &graphql.Field{Type: graphql.String},
		"order": &graphql.Field{Type: graphql.String},
		"mode":  &graphql.Field{Type: graphql.String},
	},
})

// would be different per type (e.g. People, Publications)
var sorter = graphql.NewObject(graphql.ObjectConfig{
	Name: "Sort",
	Fields: graphql.Fields{
		"fields": &graphql.Field{Type: graphql.NewList(sortField)},
	},
})

var personFacets = graphql.NewObject(graphql.ObjectConfig{
	Name: "PersonFacets",
	Fields: graphql.Fields{
		"departments": &graphql.Field{Type: graphql.NewList(facet)},
		"types":       &graphql.Field{Type: graphql.NewList(facet)},
		"keywords":    &graphql.Field{Type: graphql.NewList(facet)},
	},
})

var facet = graphql.NewObject(graphql.ObjectConfig{
	Name: "Facet",
	Fields: graphql.Fields{
		"label": &graphql.Field{Type: graphql.String},
		"count": &graphql.Field{Type: graphql.Int},
	},
})

var personList = graphql.NewObject(graphql.ObjectConfig{
	Name: "PersonList",
	Fields: graphql.Fields{
		"results":  &graphql.Field{Type: graphql.NewList(person)},
		"pageInfo": &graphql.Field{Type: pageInfo},
		"facets":   &graphql.Field{Type: personFacets},
	},
})

var grantList = graphql.NewObject(graphql.ObjectConfig{
	Name: "GrantList",
	Fields: graphql.Fields{
		"results":  &graphql.Field{Type: graphql.NewList(grant)},
		"pageInfo": &graphql.Field{Type: pageInfo},
	},
})

var publicationList = graphql.NewObject(graphql.ObjectConfig{
	Name: "publicationList",
	Fields: graphql.Fields{
		"results":  &graphql.Field{Type: graphql.NewList(publication)},
		"pageInfo": &graphql.Field{Type: pageInfo},
	},
})

// these are graphql only json models
type PageInfo struct {
	PerPage     int `json:"perPage"`
	CurrentPage int `json:"page"`
	TotalPages  int `json:"totalPages"`
	Count       int `json:"count"`
}

type Facet struct {
	Label string `json:"label"`
	Count int64  `json:"count"`
}

type PeopleFacets struct {
	Keywords    []Facet `json:"keywords"`
	Types       []Facet `json:"types"`
	Departments []Facet `json:"departments"`
}

type PersonList struct {
	Results  []Person      `json:"data"`
	PageInfo PageInfo      `json:"pageInfo"`
	Facets   *PeopleFacets `json:"facets"`
}

type PublicationList struct {
	Results  []Publication `json:"data"`
	PageInfo PageInfo      `json:"pageInfo"`
}

type GrantList struct {
	Results  []Grant  `json:"data"`
	PageInfo PageInfo `json:"pageInfo"`
}

// A way to just call queries without the server (maybe?)
// e.g.
// qry := `{
//	publicationList {
//	  results {
//		id
//		title
//	  }
//	 }
// }`
//
// result := ExecuteQuery(qry)
// or
// params := map[string]string{
//    "id": "per000001",
//}
// result := ExecuteQueryWithParams(qry, params)

// NOTE: just a wrapper - probably better to just get schema
// and then use graphql lib directly
func ExecuteQueryWithParams(query string,
	variables map[string]interface{}) *graphql.Result {
	schema := GetVivoSchema()
	result := graphql.Do(graphql.Params{
		Schema:         *schema,
		RequestString:  query,
		VariableValues: variables,
	})
	return result
}

// just in case there are no parameters
func ExecuteQuery(query string) *graphql.Result {
	schema := GetVivoSchema()
	result := graphql.Do(graphql.Params{
		Schema:        *schema,
		RequestString: query,
	})
	return result
}
