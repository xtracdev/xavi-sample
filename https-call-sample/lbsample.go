package main

import (
	"os"
	"github.com/xtracdev/xavi/env"
	"log"
	"fmt"
	"github.com/xtracdev/xavi/kvstore"
	"github.com/xtracdev/xavi/loadbalancer"
	"net/http"
	"golang.org/x/net/context"
	"io/ioutil"
	"github.com/xtracdev/xavi/config"
	_ "expvar"
)

func fatal(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {

	//Allow expvar monitoring
	go http.ListenAndServe(":1234", nil)

	//Initialize  the KVStore from the local file store established by running
	//the set up script
	endpoint := os.Getenv(env.KVStoreURL)
	if endpoint == "" {
		log.Fatal(fmt.Sprintf("Required environment variable %s for configuration KV store must be specified", env.KVStoreURL))
	}

	kvs, err := kvstore.NewKVStore(endpoint)
	fatal(err)

	//Since no listener is in context in this sample process, register
	//the configuration so NewBackendLoadBalancer has it in context.
	sc, err := config.ReadServiceConfig("quote-listener", kvs)
	fatal(err)
	config.RecordActiveConfig(sc)

	//Instantiate the load balancer for the quote-backend
	lb, err := loadbalancer.NewBackendLoadBalancer("quote-backend")
	fatal(err)

	req, err := http.NewRequest("GET", "https://hostname:4443", nil)
	fatal(err)

	for i := 1; i < 100000; i++ {
		resp, err := lb.DoWithLoadBalancer(context.Background(), req, true)
		fatal(err)

		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		fatal(err)

		fmt.Printf("Read %s\n", string(b))
	}
}
