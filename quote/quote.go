package quote

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/xtracdev/xavi/plugin"
	"github.com/xtracdev/xavisample/session"
	"golang.org/x/net/context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
)

func extractResource(uri string) (string, error) {
	parts := strings.Split(uri, "/")
	if len(parts) != 3 || parts[2] == "" {
		return "", fmt.Errorf("Expected URI format: /quote/<symbol>")
	}

	return parts[2], nil

}

func NewQuoteWrapper() plugin.Wrapper {
	return new(QuoteWrapper)
}

type QuoteWrapper struct{}

func (lw QuoteWrapper) Wrap(h plugin.ContextHandler) plugin.ContextHandler {
	return plugin.ContextHandlerFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {

		if c != nil {
			sessionId, ok := c.Value(session.SessionKey).(int)
			if ok {
				println("session id:", sessionId)
			}
		}

		//Grab the symbol to quote from the uri
		resourceId, err := extractResource(r.RequestURI)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}

		//Convert the method to POST for SOAP, and set the soap service
		//endpoint for the destination server
		r.Method = "POST"
		r.URL.Path = "/services/quote/getquote"

		//Form the SOAP payload
		payload := getQuoteRequestForSymbol(resourceId)
		payloadBytes, err := xml.Marshal(&payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//Post the payload, and record the response
		r.Body = ioutil.NopCloser(bytes.NewReader(payloadBytes))
		rec := httptest.NewRecorder()

		h.ServeHTTPContext(c, rec, r)

		//Parse the recorded response to allow the quote price to be extracted
		var response ResponseEnvelope
		err = xml.Unmarshal(rec.Body.Bytes(), &response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//Return just the price to the caller
		w.Write([]byte(response.Body.GetLastTradePriceResponse.Price + "\n"))
	})
}
