package vivographql

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/olivere/elastic"
)

func addToIndex(index string, typeName string, id string, obj interface{}) {
	ctx := context.Background()
	client := GetElasticClient()

	get1, err := client.Get().
		Index(index).
		Type(typeName).
		Id(id).
		Do(ctx)

	switch {
	case elastic.IsNotFound(err):
		put1, err := client.Index().
			Index(index).
			Type(typeName).
			Id(id).
			BodyJson(obj).
			Do(ctx)

		if err != nil {
			panic(err)
		}

		fmt.Printf("ADDED %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
		spew.Println(obj)
		return
	case elastic.IsConnErr(err):
		panic(err)
	case elastic.IsTimeout(err):
		panic(err)
	case err != nil:
		panic(err)
	}

	if get1.Found {
		update1, err := client.Update().
			RetryOnConflict(2).
			Index(index).
			Type(typeName).
			Id(id).
			Doc(obj).
			Do(ctx)

		if err != nil {
			panic(err)
		}

		fmt.Printf("UPDATED %s to index %s, type %s\n", update1.Id, update1.Index, update1.Type)
	}

	if err != nil {
		// Handle error
		panic(err)
	}
	spew.Println(obj)
}

func partialUpdate(index string, typeName string, id string, prop string, obj interface{}) {
	ctx := context.Background()
	client := GetElasticClient()

	get1, err := client.Get().
		Index(index).
		Type(typeName).
		Id(id).
		Do(ctx)

	switch {
	case elastic.IsNotFound(err):
		// NOTE: in theory we could add without source doc
		fmt.Printf("no doc id=%s found to append to\n", id)
		//fmt.Printf("ADDED %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
		return
	case elastic.IsConnErr(err):
		panic(err)
	case elastic.IsTimeout(err):
		panic(err)
	case err != nil:
		panic(err)
	}

	if get1.Found {
		update1, err := client.Update().
			RetryOnConflict(2).
			Index(index).
			Type(typeName).
			Id(id).
			//Doc(obj).
			// replace all of prop ??...
			Doc(map[string]interface{}{prop: obj}).
			DetectNoop(true).
			Do(ctx)

		if err != nil {
			panic(err)
		}

		fmt.Printf("UPDATED %s to index %s, type %s\n", update1.Id, update1.Index, update1.Type)
	}

	if err != nil {
		// Handle error
		panic(err)
	}
	spew.Println(obj)
}

func clearIndex(name string) {
	ctx := context.Background()

	client := GetElasticClient()

	deleteIndex, err := client.DeleteIndex(name).Do(ctx)
	if err != nil {
		log.Printf("ERROR:%v\n", err)
		return
	}
	if !deleteIndex.Acknowledged {
		// Not acknowledged
		log.Println("Not acknowledged")
	} else {
		log.Println("Acknowledged!")
	}
}

/*
func ClearPeopleIndex() {
	clearIndex("people")
}

func ClearAffiliationsIndex() {
	clearIndex("affiliations")
}

func ClearEducationsIndex() {
	clearIndex("educations")
}

func ClearGrantsIndex() {
	clearIndex("grants")
}

func ClearFundingRolesIndex() {
	clearIndex("funding-roles")
}

func ClearPublicationsIndex() {
	clearIndex("publications")
}

func ClearAuthorshipsIndex() {
	clearIndex("authorships")
}

*/
// make these return error?
func PersonMapping() (string, error) {
	return RenderTemplate("person.tmpl")
}

func AffiliationMapping() (string, error) {
	return RenderTemplate("affiliation.tmpl")
}

func FundingRoleMapping() (string, error) {
	return RenderTemplate("funding-role.tmpl")
}

func PublicationMapping() (string, error) {
	return RenderTemplate("publication.tmpl")
}

func AuthorshipMapping() (string, error) {
	return RenderTemplate("authorship.tmpl")
}

func GrantMapping() (string, error) {
	return RenderTemplate("grant.tmpl")
}

// NOTE: 'mappingJson' is just a json string plugged into template
func makeIndex(name string, mappingJson string) {
	ctx := context.Background()

	client := GetElasticClient()

	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists(name).Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex(name).BodyString(mappingJson).Do(ctx)
		if err != nil {
			// Handle error
			panic(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}
}

/*
func MakePeopleIndex(mapping string) {
	makeIndex("people", mapping)
}

func MakeGrantsIndex(mapping string) {
	makeIndex("grants", mapping)
}

func MakeFundingRolesIndex(mapping string) {
	makeIndex("funding-roles", mapping)
}

func MakePublicationsIndex(mapping string) {
	makeIndex("publications", mapping)
}

func MakeAuthorshipsIndex(mapping string) {
	makeIndex("authorships", mapping)
}

*/

// needs json data -->
// TODO: use Identifiable interface? e.g.
// func AddResources(resources ...Identifiable) {
//	// how to make new Person{} each loop
//}
func AddPeople(people ...string) {
	for _, element := range people {
		resource := Person{}
		data := []byte(element)
		json.Unmarshal(data, &resource)
		addToIndex("people", "person", resource.Id, resource)
	}
}

func AddAffiliationsToPeople(positions ...string) {
	// need to group by personId
	collections := make(map[string][]Affiliation)

	for _, element := range positions {
		resource := Affiliation{}
		data := []byte(element)
		json.Unmarshal(data, &resource)

		collections[resource.PersonId] = append(collections[resource.PersonId], resource)
	}

	for key, value := range collections {
		partialUpdate("people", "person", key, "affiliationList", value)
	}
}

func AddEducationsToPeople(educations ...string) {
	// need to group by personId
	collections := make(map[string][]Education)

	for _, element := range educations {
		resource := Education{}
		data := []byte(element)
		json.Unmarshal(data, &resource)

		collections[resource.PersonId] = append(collections[resource.PersonId], resource)
	}
	for key, value := range collections {
		partialUpdate("people", "person", key, "educationList", value)
	}
}

func AddGrants(grants ...string) {
	for _, element := range grants {
		resource := Grant{}
		data := []byte(element)
		json.Unmarshal(data, &resource)

		addToIndex("grants", "grant", resource.Id, resource)
	}
}

func AddFundingRoles(fundingRoles ...string) {
	for _, element := range fundingRoles {
		resource := FundingRole{}
		data := []byte(element)
		json.Unmarshal(data, &resource)

		addToIndex("funding-roles", "funding-role", resource.Id, resource)
	}
}

// need at least an id
func AddPublications(publications ...string) {
	for _, element := range publications {
		resource := Publication{}
		data := []byte(element)
		json.Unmarshal(data, &resource)

		addToIndex("publications", "publication", resource.Id, resource)
	}
}

func AddAuthorships(authorships ...string) {
	for _, element := range authorships {
		resource := Authorship{}
		data := []byte(element)
		json.Unmarshal(data, &resource)

		addToIndex("authorships", "authorship", resource.Id, resource)
	}
}
