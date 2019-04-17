package vivographql

type Config struct {
	Elastic   elasticSearch `toml:"elastic"`
	Graphql   graphqlServer `toml:"graphql"`
	Templates templatePaths `toml:"template"`
}

type elasticSearch struct {
	Url string
}

type graphqlServer struct {
	Port int
}

type templatePaths struct {
	Layout  string `toml:"layout"`
	Include string `toml:"include"`
}

type schemasPath struct {
	Dir string `toml:"dir"`
}
