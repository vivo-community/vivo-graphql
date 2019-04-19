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

	// mutations - or adding records at least
	AddPeople(people ...Person)
	// note these seem to be better was partial updates to person
	AddAffiliationsToPeople(collection map[string][]Affiliation)
	AddEducationsToPeople(collection map[string][]Affiliation)
	AddGrants(grants ...Grant)
	AddFundingRoles(fundingRoles ...FundingRole)
	AddPublications(publications ...Publication)
	AddAuthorships(authorships ...Authorship)

	// Clear functions e.g
	// ClearPeopleIndex()
	// Clear ...
}

// will be different per index engine
type Indexer interface {
	GetClient() interface{}
}

// would (long term) need some mechanism for choosing index
// engine at start up maybe:
/*
if conf.Elastic {
        if err := vq.EstablishElasticIndexer(conf.Elastic.Url); err != nil {
	        fmt.Printf("could not establish elastic client %s\n", err)
	        os.Exit(1)
	    }
    } else if conf.Solr {
        panic("solr engine non-functional")
    }
}

*/
