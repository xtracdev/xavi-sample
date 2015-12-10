//session has a simple plugin that shows how a session id could be placed in the context by a plugin
package session

import (
	"github.com/xtracdev/xavi/plugin"
	"golang.org/x/net/context"
	"math/rand"
	"net/http"
	"time"
)

type sessionKey int

const SessionKey sessionKey = 666

var seed = rand.NewSource(time.Now().UnixNano())
var gen = rand.New(seed)

type SessionWrapper struct{}

func NewSessionWrapper() plugin.Wrapper {
	return new(SessionWrapper)
}

func (lw SessionWrapper) Wrap(h plugin.ContextHandler) plugin.ContextHandler {
	return plugin.ContextHandlerFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
		if c == nil {
			c = context.Background()
		}

		sessionId := gen.Intn(999999999)

		c = context.WithValue(c, SessionKey, sessionId)

		h.ServeHTTPContext(c, w, r)

	})
}
