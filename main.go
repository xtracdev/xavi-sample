package main

import (
	"github.com/xtracdev/xavi/plugin"
	"github.com/xtracdev/xavi/runner"
	"github.com/xtracdev/xavisample/quote"
	"log"
	"os"
)

func registerPlugins() {
	err := plugin.RegisterWrapperFactory("Quote", quote.NewQuoteWrapper)
	if err != nil {
		log.Fatal("Error registering quote plugin factory")
	}
}

func main() {
	runner.Run(os.Args[1:], registerPlugins)
}
