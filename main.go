package main

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/xtracdev/xavi/config"
	"github.com/xtracdev/xavi/kvstore"
	"github.com/xtracdev/xavi/plugin"
	"github.com/xtracdev/xavi/plugin/recovery"
	"github.com/xtracdev/xavi/plugin/timing"
	"github.com/xtracdev/xavi/runner"
	"github.com/xtracdev/xavisample/quote"
	"github.com/xtracdev/xavisample/session"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func NewCustomRecoveryWrapper(args ...interface{}) plugin.Wrapper {
	return &recovery.RecoveryWrapper{
		RecoveryContext: customerRecoveryContext,
	}
}

var customerRecoveryContext = recovery.RecoveryContext{
	LogFn: func(r interface{}) {
		var err error
		switch t := r.(type) {
		case string:
			err = errors.New(t)
		case error:
			err = t
		default:
			err = errors.New("Unknown error")
		}
		log.Warn("Handled panic: ", err.Error())
	},
	ErrorMessageFn: func(r interface{}) string {
		return "Handled a panic... try again."
	},
}

func registerPlugins() {
	err := plugin.RegisterWrapperFactory("Quote", quote.NewQuoteWrapper)
	if err != nil {
		log.Fatal("Error registering quote plugin factory")
	}

	err = plugin.RegisterWrapperFactory("SessionId", session.NewSessionWrapper)
	if err != nil {
		log.Fatal("Error registering session id plugin factory")
	}

	err = plugin.RegisterWrapperFactory("Recovery", NewCustomRecoveryWrapper)
	if err != nil {
		log.Fatal("Error registering recovery plugin factory")
	}

	err = plugin.RegisterWrapperFactory("Timing", timing.NewTimingWrapper)
	if err != nil {
		log.Fatal("Error registering recovery plugin factory")
	}
}

func healthy(endpoint string, transport *http.Transport) <-chan bool {
	statusChannel := make(chan bool)

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Second,
	}

	go func() {

		log.Info("Hello there, this is a custom health check.")

		resp, err := client.Get(endpoint)
		if err != nil {
			log.Warn("Error doing get on healthcheck endpoint ", endpoint, " : ", err.Error())
			statusChannel <- false
			return
		}

		defer resp.Body.Close()
		ioutil.ReadAll(resp.Body)

		log.Infof("%s is healthy", endpoint)

		statusChannel <- resp.StatusCode == 200
	}()

	return statusChannel
}

func registerMyHealthchecks(kvs kvstore.KVStore) error {
	config.RegisterHealthCheckForBackend(kvs, "quote-backend", healthy)
	return nil
}

func main() {
	runner.AddKVSCallbackFunction(registerMyHealthchecks)
	runner.Run(os.Args[1:], registerPlugins)
}
