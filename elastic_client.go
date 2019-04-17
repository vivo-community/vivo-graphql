package vivographql

import (
	"net/http"

	"github.com/olivere/elastic"
)

//https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1

//ge "github.com/OIT-ads-web/graphql_endpoint"
//"github.com/davecgh/go-spew/spew"

var ElasticClient *elastic.Client

func GetElasticClient() *elastic.Client {
	return ElasticClient
}

func MakeElasticClient(url string) error {
	// establishing a 'global' client
	client, err := elastic.NewClient(elastic.SetURL(url),
		elastic.SetSniff(false))

	// NOTE: this is establishing a global client because the elastic client is
	// supposed to be long-lived
	// see https://github.com/olivere/elastic/blob/release-branch.v6/client.go
	ElasticClient = client
	return err
}

func MakeElasticClientDebug(url string, httpClient *http.Client) error {
	// establishing a 'global' client
	client, err := elastic.NewClient(elastic.SetURL(url),
		elastic.SetSniff(false),
		elastic.SetHttpClient(httpClient))

	// NOTE: this is establishing a global client because the elastic client is
	// supposed to be long-lived
	// see https://github.com/olivere/elastic/blob/release-branch.v6/client.go
	ElasticClient = client
	return err
}
