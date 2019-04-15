package examples

import (
	"context"
	"encoding/json"
	"log"

	el "github.com/OIT-ads-web/graphql_endpoint/elastic"
	"github.com/davecgh/go-spew/spew"
	"github.com/olivere/elastic"
)

// should eventually try stuff more like here:
// https://github.com/olivere/elastic/wiki/QueryDSL
func ExampleIdQuery() {
	el.IdQuery("people", []string{"per4774112", "per8608642"})
}

func ExampleListPeople() {
	el.ListAll("people")
}

/*
target:

aggs: {
   "types" : { "terms": {"field" : "type.label" }},
   "keywords": {
      "nested": {
          "path": "keywordList"
      },
      "aggs": {
        "keyword" : { "terms" : { "field": "keywordList.label.keyword" } }
      }
   },
}
*/
/*
getting this (kind of same):

{
   "aggregations":{
      "keywords":{
         "aggregations":{
            "keyword":{
               "terms":{
                  "field":"keywordList.label.keyword"
               }
            }
         },
         "nested":{
            "path":"keywordList"
         }
      },
      "types":{
         "terms":{
            "field":"type.label"
         }
      }
   },
   "query":{
      "match_all":{

      }
   }
}

what if:

Agg<name> -> spec ->

type Terms struct {
	Field string `json:"field"`
} `json:"terms"`

type Nested struct {
	Path string `json:"path"`
} `json:"nested"`

aggs: {
   "types" : { "terms": {"field" : "type.label" }},
   "keywords": {
      "nested": {
          "path": "keywordList"
      },
      "aggs": {
        "keyword" : { "terms" : { "field": "keywordList.label.keyword" } }
      }
   },
}

var map = map[string]string{
	"aggs": {
		"types": { "terms": {"field" : "type.label"}},
	}
}

https://play.golang.org/p/gOoyOQh9Y5

type Property interface{}

// root
type State []Property

func NState(ps ...Property) State{
	return State(ps[:])
}

var traffic = NState(
	Title("Traffic Signal"),
	Entry(RedOn),
	OnError(HandleError),
	NState(
		Title("Red"),
		Entry(RedOn),
	),
)

type Agg struct {
  terms: make(map[string][]string)
}

Aggs {
	map[string]Terms|???
}

type Aggs struct {
	map [string]Agg
} `json:"aggs"`

https://stackoverflow.com/questions/27553399/golang-how-to-initialize-a-map-field-within-a-struct

type Alpha struct {
    Pix []uint8
    Stride int
    Rect Rectangle
}

func NewAlpha(r Rectangle) *Alpha {
    w, h := r.Dx(), r.Dy()
    pix := make([]uint8, 1*w*h)
    return &Alpha{pix, 1 * w, r}
}

https://github.com/jmlucjav/elasticsearch2jsonschema

https://github.com/kristianmandrup/json-schema-to-es-mapping

https://gist.github.com/AdrianRossouw/8887766ca0e7052814b0

https://github.com/elastic/elasticsearch/tree/master/rest-api-spec/src/main/resources/rest-api-spec/api

https://grokbase.com/p/gg/elasticsearch/149t0q1tmg/loading-json-ld-into-es

*/
func ExampleAggregations() {
	ctx := context.Background()
	client := el.GetClient()

	q := elastic.NewMatchAllQuery()

	service := client.Search().
		Index("people").
		Query(q)

	agg := elastic.NewTermsAggregation().Field("type.label")
	service = service.Aggregation("types", agg)

	nested := elastic.NewNestedAggregation().Path("keywordList")
	subAgg := nested.SubAggregation("keyword",
		elastic.NewTermsAggregation().Field("keywordList.label.keyword"))

	nested2 := elastic.NewNestedAggregation().Path("affiliationList")
	subAgg2 := nested2.SubAggregation("department",
		elastic.NewTermsAggregation().Field("affiliationList.organization.label.dept"))

	service = service.Aggregation("keywords", subAgg)
	service = service.Aggregation("affiliations", subAgg2)

	searchResult, err := service.Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}

	for _, hit := range searchResult.Aggregations {
		var obj interface{}
		err := json.Unmarshal(*hit, &obj)
		if err != nil {
			panic(err)
		}
		str := spew.Sdump(obj)
		log.Println(str)
	}

	log.Println("************")
}
