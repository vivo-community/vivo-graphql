package vivographql

import (
	"net/http"

	"github.com/olivere/elastic"
)

var elasticClient *elastic.Client

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

// ??
func EstablishElasticIndexer(url string) error {
	// debug or not debug version config
	return makeElasticClient(url)
	// indexer = ElasticIndexer{ -- ?}
}

func makeElasticClient(url string) error {
	// establishing a 'global' client
	client, err := elastic.NewClient(elastic.SetURL(url),
		elastic.SetSniff(false))

	// NOTE: this is establishing a global client because the elastic client is
	// supposed to be long-lived
	// see https://github.com/olivere/elastic/blob/release-branch.v6/client.go
	elasticClient = client
	return err
}

func makeElasticClientDebug(url string, httpClient *http.Client) error {
	// establishing a 'global' client
	client, err := elastic.NewClient(elastic.SetURL(url),
		elastic.SetSniff(false),
		elastic.SetHttpClient(httpClient))

	// NOTE: this is establishing a global client because the elastic client is
	// supposed to be long-lived
	// see https://github.com/olivere/elastic/blob/release-branch.v6/client.go
	elasticClient = client
	return err
}
