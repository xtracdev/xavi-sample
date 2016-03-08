package main

import (
	"github.com/xtracdev/xavi/plugin"
	"github.com/xtracdev/xavi/runner"
	"github.com/xtracdev/xavisample/quote"
	"github.com/xtracdev/xavisample/session"
	log "github.com/Sirupsen/logrus"
	"os"
	"github.com/xtracdev/xavi/plugin/recovery"
	"errors"
	"github.com/xtracdev/xavi/plugin/timing"
)

func NewCustomRecoveryWrapper() plugin.Wrapper {
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

func main() {
	runner.Run(os.Args[1:], registerPlugins)
}
