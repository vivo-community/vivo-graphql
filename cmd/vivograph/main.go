package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/rs/cors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	vq "github.com/vivo-community/vivo-graphql"
)

var conf vq.Config

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

	vq.LoadSchemas(conf)

	c := cors.New(cors.Options{
		AllowCredentials: true,
	})

	handler := vq.MakeHandler()
	http.Handle("/graphql", c.Handler(handler))

	port := conf.Graphql.Port
	portConfig := fmt.Sprintf(":%d", port)
	err := http.ListenAndServe(portConfig, nil)
	if err != nil {
		fmt.Printf("server start error: %v\n", err)
	}

}
