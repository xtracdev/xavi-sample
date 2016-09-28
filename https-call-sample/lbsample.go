package main

import (
	_ "expvar"
	"fmt"
	"github.com/xtracdev/xavi/config"
	"github.com/xtracdev/xavi/env"
	"github.com/xtracdev/xavi/kvstore"
	"github.com/xtracdev/xavi/loadbalancer"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
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

	for i := 1; i < 10000; i++ {
		req, err := http.NewRequest("GET", "https://hostname:4443", nil)
		fatal(err)

		resp, err := lb.DoWithLoadBalancer(req, true)
		fatal(err)

		_, err = ioutil.ReadAll(resp.Body)
		fatal(err)
		resp.Body.Close()

		if i%100 == 0 {
			fmt.Printf("Done %d calls...\n", i)
		}
	}

	time.Sleep(300 * time.Second)

}
