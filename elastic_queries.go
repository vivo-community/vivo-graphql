package vivographql

import (
	"context"
	"encoding/json"
	"log"

	"github.com/olivere/elastic"
)

func (indexer *ElasticIndexer) FindPerson(personId string) (Person, error) {
	var person = Person{}

	ctx := context.Background()
	client := indexer.GetClient()

	log.Printf("looking for person %s\n", personId)

	get1, err := client.Get().
		Index("people").
		Id(personId).
		Do(ctx)
	if err != nil {
		return person, err
	}

	err = json.Unmarshal(*get1.Source, &person)
	return person, err
}

func (indexer *ElasticIndexer) FindPublication(publicationId string) (Publication, error) {
	var publication = Publication{}

	ctx := context.Background()
	client := indexer.GetClient()

	log.Printf("looking for publication %s\n", publicationId)

	get1, err := client.Get().
		Index("publications").
		Id(publicationId).
		Do(ctx)
	if err != nil {
		return publication, err
	}

	err = json.Unmarshal(*get1.Source, &publication)
	return publication, err
}

func (indexer *ElasticIndexer) FindGrant(grantId string) (Grant, error) {
	var grant = Grant{}

	ctx := context.Background()
	client := indexer.GetClient()

	log.Printf("looking for grant %s\n", grantId)

	get1, err := client.Get().
		Index("grants").
		Id(grantId).
		Do(ctx)
	if err != nil {
		return grant, err
	}

	err = json.Unmarshal(*get1.Source, &grant)
	return grant, err
}

func parsePeopleAggregations(facets elastic.Aggregations) *PeopleFacets {
	peopleFacets := &PeopleFacets{}

	if agg, found := facets.Nested("keywords"); found {
		var facets []Facet
		if sub, subFound := agg.Terms("keyword"); subFound {
			for _, bucket := range sub.Buckets {
				facet := Facet{Label: bucket.Key.(string), Count: bucket.DocCount}
				facets = append(facets, facet)
			}
		}
		peopleFacets.Keywords = facets
	}

	if agg, found := facets.Nested("affiliations"); found {
		var facets []Facet
		if sub, subFound := agg.Terms("department"); subFound {
			for _, bucket := range sub.Buckets {
				facet := Facet{Label: bucket.Key.(string), Count: bucket.DocCount}
				facets = append(facets, facet)
			}
		}
		peopleFacets.Departments = facets
	}

	if agg, found := facets.Terms("types"); found {
		var facets []Facet
		for _, bucket := range agg.Buckets {
			facet := Facet{Label: bucket.Key.(string), Count: bucket.DocCount}
			facets = append(facets, facet)
		}
		peopleFacets.Types = facets
	}
	return peopleFacets
}

func (indexer *ElasticIndexer) FindPeople(limit int, offset int, query string) (PersonList, error) {
	var people []Person
	ctx := context.Background()
	client := indexer.GetClient()

	q := elastic.NewQueryStringQuery(query)
	log.Println("looking for people")

	service := client.Search().
		Index("people").
		Query(q).
		From(offset).
		Size(limit)

	keywordsSize := 100
	departmentsSize := 100

	// TODO: kind of kludged together here - probably much better way to do
	agg := elastic.NewTermsAggregation().Field("type.label")
	service = service.Aggregation("types", agg)

	nested := elastic.NewNestedAggregation().Path("keywordList")
	subAgg := nested.SubAggregation("keyword",
		elastic.NewTermsAggregation().
			Field("keywordList.label.keyword").
			Size(keywordsSize))

	nested2 := elastic.NewNestedAggregation().Path("affiliationList")
	subAgg2 := nested2.SubAggregation("department",
		elastic.NewTermsAggregation().
			Field("affiliationList.organization.label.dept").
			Size(departmentsSize))

	service = service.Aggregation("keywords", subAgg)
	service = service.Aggregation("affiliations", subAgg2)

	searchResult, err := service.Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}

	for _, hit := range searchResult.Hits.Hits {
		person := Person{}
		err := json.Unmarshal(*hit.Source, &person)
		if err != nil {
			panic(err)
		}
		people = append(people, person)
	}

	// TODO: might be one off
	totalHits := int(searchResult.TotalHits())
	log.Printf("total people hits: %d\n", totalHits)

	pageInfo := FigurePaging(limit, offset, totalHits)
	facets := parsePeopleAggregations(searchResult.Aggregations)
	personList := PersonList{Results: people, PageInfo: pageInfo, Facets: facets}
	return personList, err
}

func (indexer *ElasticIndexer) FindPublications(limit int, offset int, query string) (PublicationList, error) {
	var publications []Publication
	ctx := context.Background()
	// should query elastic here
	client := indexer.GetClient()

	//q := elastic.NewMatchAllQuery()
	q := elastic.NewQueryStringQuery(query)
	service := client.Search().
		Index("publications").
		Query(q).
		From(offset).
		Size(limit)

	/*
		// TODO: kind of kludged together here - probably much better way to do
		agg := elastic.NewTermsAggregation().Field("type.label")
		service = service.Aggregation("types", agg)

		nested := elastic.NewNestedAggregation().Path("keywordList")
		subAgg := nested.SubAggregation("keyword", elastic.NewTermsAggregation().Field("keywordList.label.keyword"))

		nested2 := elastic.NewNestedAggregation().Path("affiliationList")
		subAgg2 := nested2.SubAggregation("department",
			elastic.NewTermsAggregation().Field("affiliationList.organization.label.dept"))

		service = service.Aggregation("keywords", subAgg)
		service = service.Aggregation("affiliations", subAgg2)
	*/

	searchResult, err := service.Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}

	for _, hit := range searchResult.Hits.Hits {
		publication := Publication{}
		err := json.Unmarshal(*hit.Source, &publication)
		if err != nil {
			panic(err)
		}
		publications = append(publications, publication)
	}

	// might be one off
	totalHits := int(searchResult.TotalHits())
	log.Printf("total publication hits: %d\n", totalHits)

	pageInfo := FigurePaging(limit, offset, totalHits)
	// eventually
	//facets := parsePublicationsAggregations(searchResult.Aggregations)
	//publicationList := PublicationList{Results: publications, PageInfo: pageInfo, Facets: facets}
	publicationList := PublicationList{Results: publications, PageInfo: pageInfo}
	return publicationList, err
	//return publications, err
}

func (indexer *ElasticIndexer) FindPersonPublications(personId string, limit int, offset int) (PublicationList, error) {
	var publications []Publication
	var publicationIds []string

	ctx := context.Background()
	client := indexer.GetClient()

	q := elastic.NewMatchQuery("personId", personId)

	searchResult, err := client.Search().
		Index("authorships").
		Query(q).
		From(offset).
		Size(limit).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}

	// FIXME: could optimize better - dataloader etc...
	for _, hit := range searchResult.Hits.Hits {
		authorship := Authorship{}
		err := json.Unmarshal(*hit.Source, &authorship)
		if err != nil {
			panic(err)
		}

		publicationId := authorship.PublicationId
		publicationIds = append(publicationIds, publicationId)
	}

	// NOTE: need to have the count be authorship search
	// not publication search - since pub search is just
	// an id search derived from authorship search

	// NOTE: not sure this is actually correct, might be one off
	totalHits := int(searchResult.TotalHits())
	log.Printf("total authorships: %d\n", totalHits)

	// ids query
	pubQuery := elastic.NewIdsQuery("publication").
		Ids(publicationIds...)

	pubResults, err := client.Search().
		Index("publications").
		Query(pubQuery).
		From(0).
		Size(totalHits).
		RequestCache(true).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}

	for _, hit := range pubResults.Hits.Hits {
		publication := Publication{}
		err := json.Unmarshal(*hit.Source, &publication)
		if err != nil {
			panic(err)
		}
		publications = append(publications, publication)
	}

	log.Printf("size: %d, from:%d\n", limit, offset)

	pageInfo := FigurePaging(limit, offset, totalHits)
	publicationList := PublicationList{Results: publications, PageInfo: pageInfo}

	return publicationList, err
}

func (indexer *ElasticIndexer) FindGrants(limit int, offset int, query string) (GrantList, error) {
	var grants []Grant
	ctx := context.Background()
	client := indexer.GetClient()

	//q := elastic.NewMatchAllQuery()
	q := elastic.NewQueryStringQuery(query)
	service := client.Search().
		Index("grants").
		Query(q).
		From(offset).
		Size(limit)

	searchResult, err := service.Do(ctx)

	if err != nil {
		// Handle error
		panic(err)
	}

	for _, hit := range searchResult.Hits.Hits {
		grant := Grant{}
		err := json.Unmarshal(*hit.Source, &grant)
		if err != nil {
			panic(err)
		}
		grants = append(grants, grant)
	}

	// is this the correct number?
	totalHits := int(searchResult.TotalHits())
	log.Printf("total grant hits: %d\n", totalHits)

	pageInfo := FigurePaging(limit, offset, totalHits)
	// eventually
	//facets := parseGrantsAggregations(searchResult.Aggregations)
	//grantList := GrantList{Results: grants, PageInfo: pageInfo, Facets: facets}
	grantList := GrantList{Results: grants, PageInfo: pageInfo}
	return grantList, err
}

func (indexer *ElasticIndexer) FindPersonGrants(personId string, limit int, offset int) (GrantList, error) {
	var grants []Grant
	var grantIds []string

	ctx := context.Background()
	client := indexer.GetClient()

	q := elastic.NewMatchQuery("personId", personId)

	searchResult, err := client.Search().
		Index("funding-roles").
		Query(q).
		From(offset).
		Size(limit).
		Do(ctx)
	if err != nil {
		// handle error
		panic(err)
	}

	// is this the correct number?
	totalHits := int(searchResult.TotalHits())
	log.Printf("total funding-roles: %d\n", totalHits)

	// fixme: could optimize better - dataloader etc...
	for _, hit := range searchResult.Hits.Hits {
		fundingRole := FundingRole{}
		err := json.Unmarshal(*hit.Source, &fundingRole)
		if err != nil {
			panic(err)
		}

		grantId := fundingRole.GrantId
		grantIds = append(grantIds, grantId)
	}

	grantQuery := elastic.NewIdsQuery("grant").
		Ids(grantIds...)

	grantResults, err := client.Search().
		Index("grants").
		Query(grantQuery).
		From(0).
		Size(totalHits).
		RequestCache(true).
		Do(ctx)
	if err != nil {
		// handle error
		panic(err)
	}
	for _, hit := range grantResults.Hits.Hits {
		grant := Grant{}
		err := json.Unmarshal(*hit.Source, &grant)
		if err != nil {
			panic(err)
		}
		grants = append(grants, grant)
	}

	pageInfo := FigurePaging(limit, offset, totalHits)
	grantList := GrantList{Results: grants, PageInfo: pageInfo}
	return grantList, err
}

// remaining are just debug/util functions
/*
func ListAll(index string) {
	ctx := context.Background()
	client := GetElasticClient()
	q := elastic.NewMatchAllQuery()

	searchResult, err := client.Search().
		Index(index).
		//Type().
		Query(q).
		From(100).
		Size(100).
		Pretty(true).
		// Timeout("1000ms"). or
		// Timeout(1000).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}

	fmt.Println("********* BEGIN **********")
	for _, hit := range searchResult.Hits.Hits {
		var obj interface{}
		err := json.Unmarshal(*hit.Source, &obj)
		if err != nil {
			panic(err)
		}
		spew.Printf("%v\n", obj)
	}
	fmt.Printf("********* END (%d) **********\n", searchResult.TotalHits())
}

func IdQuery(index string, ids []string) {
	// NOTE: can send 'type' into NewIdsQuery
	q := elastic.NewIdsQuery().Ids(ids...) //.QueryName("my_query")
	ctx := context.Background()
	client := GetElasticClient()

	searchResult, err := client.Search().
		Index(index).
		//Type().
		Query(q).
		From(0).
		Size(1000).
		Pretty(true).
		// Timeout("1000ms"). or
		// Timeout(1000).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}

	fmt.Println("********* BEGIN **********")
	for _, hit := range searchResult.Hits.Hits {
		var obj interface{}
		err := json.Unmarshal(*hit.Source, &obj)
		if err != nil {
			panic(err)
		}
		spew.Printf("%v\n", obj)
	}
	fmt.Println("************** END **********")
}

func FindOne(index string, personId string) {
	ctx := context.Background()
	client := GetElasticClient()

	get1, err := client.Get().
		Index(index).
		Id(personId).
		Do(ctx)

	switch {
	case elastic.IsNotFound(err):
		fmt.Println("404 not found")
	case elastic.IsConnErr(err):
		fmt.Println("connectino error")
	case elastic.IsTimeout(err):
		fmt.Println("timeout")
	case err != nil:
		panic(err)
	}

	var obj interface{}
	err = json.Unmarshal(*get1.Source, &obj)
	if err != nil {
		panic(err)
	}
	spew.Printf("%v\n", obj)
}
*/
