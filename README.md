# graphql server in golang

## config 

can either set environmental variables, 

> `export ELASTIC_URL="http://localhost:9200"`
> `export GRAPHQL_PORT="9001"`
> `export MAPPING_TEMPLATES_LAYOUT="schema-templates/layout/"`
> `export MAPPING_TEMPLATES_INCLUDE="schema-templates/"`
> `export SCHEMAS_DIR="json-schemas/"`

or if `set ENVIRONMENT=development` looks for config.toml file
in current directory (see config.toml.example)

## server 

* endpoint on `GRAPHQL_PORT`
* see localhost:<GRAPHQL_PORT>/graphql
