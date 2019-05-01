package vivographql

type Config struct {
	Elastic ElasticSearch `toml:"elastic"`
	Graphql GraphqlServer `toml:"graphql"`
	Mapping Mapping       `toml:"mapping"`
	Schemas SchemasPath   `toml:"schema"`
}

type ElasticSearch struct {
	Url string
	// TestUrl (url for tests ...???)
}

type GraphqlServer struct {
	Port int
}

type Mapping struct {
	Templates TemplatePaths `toml:"templates"`
}

type TemplatePaths struct {
	Layout  string `toml:"layout"`
	Include string `toml:"include"`
}

type SchemasPath struct {
	Dir string `toml:"dir"`
}
