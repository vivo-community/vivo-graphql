package vivographql

type Config struct {
	Elastic elasticSearch `toml:"elastic"`
	Graphql graphqlServer `toml:"graphql"`
	Mapping mapping       `toml:"mapping"`
	Schemas schemasPath   `toml:"schema"`
}

type elasticSearch struct {
	Url string
	// TestUrl (url for tests ...???)
}

type graphqlServer struct {
	Port int
}

type mapping struct {
	Templates templatePaths `toml:"templates"`
	//Layout  string `toml:"templates_layout"`
	//Include string `toml:"templates_include"`
}

type templatePaths struct {
	Layout  string `toml:"layout"`
	Include string `toml:"include"`
}

type schemasPath struct {
	Dir string `toml:"dir"`
}
