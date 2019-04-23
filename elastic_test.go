package vivographql_test

import (
	"log"
	"os"
	"testing"

	vq "github.com/vivo-community/vivo-graphql"
)

// TODO: add some reference data here?
// func TestMain() {
//  need to make use of docker-compose-test.yml
//}
func TestElasticIdQuery(t *testing.T) {
	// TODO: how to have test data but also
	// persisted volume for dev ????
	// always specific value (instead of config)
	url := "http://localhost:9200" // different port?
	if err := vq.EstablishElasticIndexer(url); err != nil {
		log.Printf("no elastic test instance %s\n", err)
		os.Exit(1)
	}

	defer vq.GetElasticIndexer().Stop()

	// TODO: first need to add people - than check if exist etc...
}
