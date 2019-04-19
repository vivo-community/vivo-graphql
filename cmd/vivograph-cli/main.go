package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	vq "github.com/vivo-community/vivo-graphql"
)

var conf vq.Config

func outputMappings() {
	// TODO: this is only one, and doesn't allow you to point
	// to where they are generated
	mapping, err := vq.PersonMapping()
	if err != nil {
		fmt.Println(err)
	}

	err = ioutil.WriteFile("person-mapping.json", []byte(mapping), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.SetOutput(os.Stdout)

	viper.SetDefault("elastic.url", "http://localhost:9200")
	viper.SetDefault("graphql.port", "9001")

	// way to add more paths ??
	if os.Getenv("ENVIRONMENT") == "development" {
		viper.SetConfigName("config")
		viper.SetConfigType("toml")
		viper.AddConfigPath(".")

		value, exists := os.LookupEnv("CONFIG_PATH")
		if exists {
			viper.AddConfigPath(value)
		}
		viper.ReadInConfig()

	} else {
		replacer := strings.NewReplacer(".", "_")
		viper.SetEnvKeyReplacer(replacer)

		// TODO: check for error?
		viper.BindEnv("elastic.url")
		viper.BindEnv("graphql.port")
	}

	mappings := flag.Bool("mappings", false, "whether to generate elastic mapping files")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	// check for err?
	viper.BindPFlags(pflag.CommandLine)

	if err := viper.Unmarshal(&conf); err != nil {
		fmt.Printf("could not establish read into conf structure %s\n", err)
		os.Exit(1)
	}

	// NOTE: just doing elastic right now
	if err := vq.EstablishElasticIndexer(conf.Elastic.Url); err != nil {
		fmt.Printf("could not establish elastic client %s\n", err)
		os.Exit(1)
	}

	vq.LoadJsonSchemas(conf)
	// CLI for validating - or just new cmd ??

	// TODO: should say which mappings, folder to output to etc...
	if *mappings {
		vq.LoadMappingTemplates(conf)
		outputMappings()
	}

}
