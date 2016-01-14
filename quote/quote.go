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
	"github.com/xtracdev/xavi/plugin/timing"
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

		//Extract the timer from the service context
		end2endTimer := timing.TimerFromContext(c)
		if end2endTimer == nil {
			http.Error(w, "No timer in call context", http.StatusInternalServerError)
			return
		}

		//Set the top level name we want to use for recording timings, counts, etc.
		end2endTimer.Name = "Quote"

		contributor := end2endTimer.StartContributor("quote svc plugin")

		//Grab the symbol to quote from the uri
		resourceId, err := extractResource(r.RequestURI)
		if err != nil {
			contributor.End(err)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}

		if c != nil {
			sid, ok := c.Value(session.SessionKey).(int)
			if ok {
				println("-----> session:", sid, "symbol", resourceId)
			}
		}

		//Convert the method to POST for SOAP, and set the soap service
		//endpoint for the destination server
		r.Method = "POST"
		r.URL.Path = "/services/quote/getquote"

		//Form the SOAP payload
		payload := getQuoteRequestForSymbol(resourceId)
		payloadBytes, err := xml.Marshal(&payload)
		if err != nil {
			contributor.End(err)
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
			contributor.End(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//Return just the price to the caller
		w.Write([]byte(response.Body.GetLastTradePriceResponse.Price + "\n"))
		contributor.End(nil)
	})
}
