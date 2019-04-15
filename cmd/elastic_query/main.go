package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	ge "github.com/OIT-ads-web/graphql_endpoint"
	"github.com/OIT-ads-web/graphql_endpoint/elastic"
	"github.com/OIT-ads-web/graphql_endpoint/examples"
	"github.com/OIT-ads-web/graphql_endpoint/http"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func example1() {
	examples.ExampleAggregations()
}

// just a few simple functions to print out data
func main() {
	var conf ge.Config
	start := time.Now()

	viper.SetDefault("elastic.url", "http://localhost:9200")

	if os.Getenv("ENVIRONMENT") == "development" {
		viper.SetConfigName("config")
		viper.SetConfigType("toml")
		viper.AddConfigPath(".")

		value, exists := os.LookupEnv("CONFIG_PATH")
		if exists {
			viper.AddConfigPath(value)
		}

		if err := viper.ReadInConfig(); err != nil {
			fmt.Printf("could not read command line config %s\n", err)
			os.Exit(1)
		}

	} else {
		replacer := strings.NewReplacer(".", "_")
		viper.SetEnvKeyReplacer(replacer)
		if err := viper.BindEnv("elastic.url"); err != nil {
			fmt.Printf("could not read command line flags %s\n", err)
			os.Exit(1)
		}
	}

	if err := viper.Unmarshal(&conf); err != nil {
		fmt.Printf("could not establish read into conf structure %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("trying to connect to elastic at %s\n", conf.Elastic.Url)

	httpClient := http.LoggingClient

	if err := elastic.MakeClientDebug(conf.Elastic.Url, httpClient); err != nil {
		fmt.Printf("could not establish elastic client %s\n", err)
		os.Exit(1)
	}

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		fmt.Printf("could not read command line flags %s\n", err)
		os.Exit(1)
	}

	fmt.Println("******* aggregations ****")
	// cmd switch for example?
	example1()

	defer elastic.Client.Stop()

	elapsed := time.Since(start)
	fmt.Printf("%s\n", elapsed)
}
