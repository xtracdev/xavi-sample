//session implements a context aware plugin that can add a session id
package session

import (
	"context"
	"fmt"
	"github.com/xtracdev/xavi/config"
	"github.com/xtracdev/xavi/plugin"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type sessionKey int

const SessionKey sessionKey = 111

func NewSessionWrapper(args ...interface{}) plugin.Wrapper {
	fmt.Printf("Active listener names%v\n", config.ActiveListenerNames())
	for _, al := range config.ActiveListenerNames() {
		sc := config.ActiveConfigForListener(al)
		switch sc {
		case nil:
			fmt.Println("Nil service config for", al)
		default:
			sc.LogConfig()
		}
	}
	return new(SessionWrapper)
}

var seed = rand.NewSource(time.Now().UnixNano())
var gen = rand.New(seed)
var mutex sync.Mutex

type SessionWrapper struct{}

func (lw SessionWrapper) Wrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		mutex.Lock()
		val := gen.Intn(999999999)
		mutex.Unlock()

		newR := r.WithContext(context.WithValue(r.Context(), SessionKey, val))

		h.ServeHTTP(w, newR)
	})
}
