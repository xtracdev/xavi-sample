//session implements a context aware plugin that can add a session id
package session

import (
	"fmt"
	"github.com/xtracdev/xavi/config"
	"github.com/xtracdev/xavi/plugin"
	"golang.org/x/net/context"
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

func (lw SessionWrapper) Wrap(h plugin.ContextHandler) plugin.ContextHandler {
	return plugin.ContextHandlerFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {

		if c == nil {
			c = context.Background()
		}

		mutex.Lock()
		val := gen.Intn(999999999)
		mutex.Unlock()
		c = context.WithValue(c, SessionKey, val)

		h.ServeHTTPContext(c, w, r)
	})
}
