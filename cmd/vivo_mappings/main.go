package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
	vq "github.com/vivo-community/vivo-graphql"
)

var conf vq.Config

func outputSchemas() {

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

	//viper.SetDefault("elastic.url", "http://localhost:9200")
	//viper.SetDefault("graphql.port", "9001")

	// way to add more paths ??
	if os.Getenv("ENVIRONMENT") == "development" {
		viper.SetConfigName("config")
		viper.SetConfigType("toml")
		viper.AddConfigPath(".")

		value, exists := os.LookupEnv("CONFIG_PATH")
		if exists {
			viper.AddConfigPath(value)
		}
		err := viper.ReadInConfig()
		if err != nil {
			fmt.Printf("could not read config: %v\n", err)
		}
	} else {
		replacer := strings.NewReplacer(".", "_")
		viper.SetEnvKeyReplacer(replacer)

		// don't actually need these here
		viper.BindEnv("elastic.url")
		viper.BindEnv("graphql.port")

		viper.BindEnv("mapping.templates.layout")
		viper.BindEnv("mapping.templates.include")
	}

	if err := viper.Unmarshal(&conf); err != nil {
		fmt.Printf("could not establish read into conf structure %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("ELASTIC_URL=%s\n", conf.Elastic.Url)
	fmt.Printf("GRAPHQL_PORT=%d\n", conf.Graphql.Port)

	fmt.Printf("MAPPING_TEMPLATES_LAYOUT=%s\n", conf.Mapping.Templates.Layout)
	fmt.Printf("MAPPING_TEMPLATES_INCLUDE=%s\n", conf.Mapping.Templates.Include)
	/*
		if err := vq.MakeElasticClient(conf.Elastic.Url); err != nil {
			fmt.Printf("could not establish elastic client %s\n", err)
			os.Exit(1)
		}
	*/

	vq.LoadTemplates(conf)

	outputSchemas()

}
