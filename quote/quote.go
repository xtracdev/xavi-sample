package quote

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/xtracdev/xavi/plugin"
	"github.com/xtracdev/xavi/plugin/timing"
	"github.com/xtracdev/xavisample/session"
	"golang.org/x/net/context"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"time"
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

//For use in generating a variety of service names for use in exploring log management solutions,
//we'll randomly generate a service name.

var serviceNames = []string{"alpha", "bravo", "charlie", "delta", "echo", "foxtrot", "golf", "india", "hotel"}

func generateServiceName() string {
	return serviceNames[rand.Intn(len(serviceNames))]
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
		end2endTimer.Name = fmt.Sprintf("%s-quote", generateServiceName())

		contributor := end2endTimer.StartContributor("QuoteSvc.GetTradePrice")

		//Grab the symbol to quote from the uri
		resourceId, err := extractResource(r.RequestURI)
		if err != nil {
			contributor.End(err)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}

		if resourceId == "XTRAC" {
			panic(errors.New("priceless!"))
		}

		if c != nil {
			sid, ok := c.Value(session.SessionKey).(int)
			if ok {
				log.Println("session:", sid, "symbol", resourceId)
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

		c = timing.AddServiceNameToContext(c, "QuoteSoapService")

		c, cf := context.WithCancel(c)

		maybeTimeout := os.Getenv("MAYBE_TIMEOUT") != "" && rand.Float64() > 0.75

		var wg sync.WaitGroup
		wg.Add(1)

		if maybeTimeout {
			go func(ctx context.Context) {
				defer wg.Done()
				ctx, cf := context.WithTimeout(ctx, 100*time.Millisecond)
				defer cf()
				h.ServeHTTPContext(ctx, rec, r)
			}(c)
		} else {
			go func() {
				defer wg.Done()
				h.ServeHTTPContext(c, rec, r)
			}()
		}

		maybeCancel(cf)

		wg.Wait()

		//Throw in a random service delay
		delay := rand.Intn(100) + 1
		time.Sleep(time.Duration(delay) * time.Millisecond)

		//Was there an error?
		if rec.Code > 299 {
			log.Println("SOAP service returned error")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(rec.Body.Bytes())
			contributor.End(err)
			return
		}

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

func maybeCancel(cf context.CancelFunc) {
	if os.Getenv("MAYBE_CANCEL") == "" {
		return
	}

	//Flip a coin
	if rand.Float64() > 0.75 {
		log.Println("coin toss says cancel")
		cf()
	}
}
