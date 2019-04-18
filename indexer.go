package vivographql

// idea is to pave way for either solr or elastic - even though elastic
// is the only one right now
type IndexEngine interface {
	FindPerson(id string) (Person, error)
	FindPublication(publicationId string) (Publication, error)
	FindGrant(grantId string) (Grant, error)
	FindPeople(limit int, offset int, query string) (PersonList, error)
	FindPublications(limit int, offset int, query string) (PublicationList, error)
	FindPersonPublications(personId string, limit int, offset int) (PublicationList, error)
	FindGrants(limit int, offset int, query string) (GrantList, error)
	FindPersonGrants(personId string, limit int, offset int) (GrantList, error)

	//AddPeople(people ...Person) error or errors or validation results ...
}

// will be different per index engine
type Indexer interface {
	GetClient() interface{}
}

// if conf.Elastic {
/*
if err := vq.EstablishElasticIndexer(conf.Elastic.Url); err != nil {
	fmt.Printf("could not establish elastic client %s\n", err)
	os.Exit(1)
} else if conf.Solr {
panic("solr engine non-functional")
}
*/
