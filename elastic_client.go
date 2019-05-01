package vivographql

import (
	"net/http"
	"sync"

	"github.com/olivere/elastic"
)

var elasticClient *elastic.Client

var oneElastic sync.Once

func getElasticClient() *elastic.Client {
	return elasticClient
}

var elasticIndexer *ElasticIndexer

type ElasticIndexer struct {
	/// empty ??
}

func GetElasticIndexer() *ElasticIndexer {
	return elasticIndexer
}

func (indexer *ElasticIndexer) GetName() interface{} {
	return "Elastic"
}

func (indexer *ElasticIndexer) GetClient() *elastic.Client {
	return getElasticClient()
}

func (indexer *ElasticIndexer) Stop() {
	getElasticClient().Stop()
}

// ??
func EstablishElasticIndexer(url string) error {
	// debug or not debug version config
	return makeElasticClient(url)
}

func makeElasticClient(url string) error {
	// establishing a 'global' client
	var err error
	oneElastic.Do(func() {
		client, eErr := elastic.NewClient(elastic.SetURL(url),
			elastic.SetSniff(false))

		// NOTE: this is establishing a global client because the elastic client is
		// supposed to be long-lived
		// see https://github.com/olivere/elastic/blob/release-branch.v6/client.go
		elasticClient = client
		err = eErr
	})
	return err
}

// FIXME: need some flag to switch to this maybe
func makeElasticClientDebug(url string, httpClient *http.Client) error {
	var err error
	oneElastic.Do(func() {
		client, eErr := elastic.NewClient(elastic.SetURL(url),
			elastic.SetSniff(false),
			elastic.SetHttpClient(httpClient))

		elasticClient = client
		err = eErr
	})
	return err
}
