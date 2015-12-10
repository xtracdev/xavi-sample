package main

import (
	"github.com/xtracdev/xavi/plugin"
	"github.com/xtracdev/xavi/runner"
	"github.com/xtracdev/xavisample/quote"
	"github.com/xtracdev/xavisample/session"
	"log"
	"os"
)

func registerPlugins() {
	err := plugin.RegisterWrapperFactory("Quote", quote.NewQuoteWrapper)
	if err != nil {
		log.Fatal("Error registering quote plugin factory")
	}

	err = plugin.RegisterWrapperFactory("SessionId", session.NewSessionWrapper)
	if err != nil {
		log.Fatal("Error registering session id plugin factory")
	}
}

func main() {
	runner.Run(os.Args[1:], registerPlugins)
}
