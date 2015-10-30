## XAVI Sample

This provides a sample that shows how to provide a RESTful facade on top
of a soap service using [Xavi](https://github.com/xtracdev/xavi).

### Service Imposter

To simulate a soap service, we'll use [mountebank](http://www.mbtest.org/)

The imposters.json file contains a simple definition to simulate the stock
quote services example from the [SOAP 1.1](http://www.w3.org/TR/2000/NOTE-SOAP-20000508/) spec:

```xml
{
  "port": 4545,
  "protocol": "http",
  "stubs": [
    {
      "responses": [
        {
          "is": {
            "statusCode": 200,
            "body": "<SOAP-ENV:Envelope
  xmlns:SOAP-ENV=\"http://schemas.xmlsoap.org/soap/envelope/\"
  SOAP-ENV:encodingStyle=\"http://schemas.xmlsoap.org/soap/encoding/\">
   <SOAP-ENV:Body>
       <m:GetLastTradePriceResponse xmlns:m=\"Some-URI\">
           <Price>34.5</Price>
       </m:GetLastTradePriceResponse>
   </SOAP-ENV:Body>
</SOAP-ENV:Envelope>\n"
          }
        }
      ],
      "predicates": [
              {
                "and": [
                  {"equals": {"path": "/services/quote/getquote","method": "POST"}},
                  { "contains": { "body": "Envelope" } }
                ]
              }
            ]
    }
  ]
}
```

Start mountebank, and configure it to expose the above endpoint:

<pre>
curl -i -X POST -H 'Content-Type: application/json' -d@imposter.json http://127.0.0.1:2525/imposters
</pre>

Based on the above, if a request is posted to the /services/quotes/getquote endpoint,
and the body contains 'Envelope', a SOAP response is returned.

What if we wanted a simple RESTful endpoint, where we could obtain a quote
for a symbol via a simple get on resource like /quote/symbol?

Doing this with Xavi is simple - we write a plug for the framework, register
it on startup, and provide some configuration with a route .

### Writing the plugin

For this example, you need to have golang 1.5.x and enable vendoring support
by setting the GO15VENDOREXPERIMENT environment variable to 1.

First, we grab Xavi from github via go get.

<pre>
go get github.com/xtracdev/xavi
</pre>

Next, we create the plugin by creating a wrapper type, a method to instantiate the wrapper, and
implement the Wrap method.  

The details are in the quote package.

### Registering the Plugin

The plugin will be referenced by name in the Xavi route configuration. To make
the plugin available to the configuration interface, it needs to be
registered when the application is started.

Xavi provides a package that serves as an entry point to applications built using
the toolkit. Applications pass the command line arguments and a function to
register the plugins.

Refer to main.go for details.

### Xavi configuration

Once the plugin and main function are available, build the application, and
use it to configure the servers, backends, routes, and listener. The below
configuration assumes the sample was built using `go build`.

Note that before you run the commands below you need to set the
value of the XAVI_KVSTORE_URL environment variable, e.g.

<pre>
export XAVI_KVSTORE_URL=file:///`pwd`/config
</pre>

<pre>
./xavisample add-server -address localhost -port 4545 -name quotesvr1
./xavisample add-backend -name quote-backend -servers quotesvr1
./xavisample add-route -name quote-route -backend quote-backend -base-uri /quote/ -plugins Quote
./xavisample add-listener -name quote-listener -routes quote-route
</pre>

At this point, the listener can be started.

Boot the listener:

<pre>
./xavisample listen -ln quote-listener -address 0.0.0.0:8080
</pre>

The xavi endpoint is now ready to go.

### Trying it out

First, you can directly access the SOAP endpoint:


```xml
curl -X POST -d '
<SOAP-ENV:Envelope
  xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/"
  SOAP-ENV:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
   <SOAP-ENV:Body>
       <m:GetLastTradePrice xmlns:m="Some-URI">
           <symbol>DIS</symbol>
       </m:GetLastTradePrice>
   </SOAP-ENV:Body>
</SOAP-ENV:Envelope>
' localhost:4545/services/quote/getquote
```

This should produce the following response:

```xml
<SOAP-ENV:Envelope  
  xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/"  
  SOAP-ENV:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/
">   
  <SOAP-ENV:Body>
    <m:GetLastTradePriceResponse xmlns:m="Some-URI">           
      <Price>34.5</Price>
    </m:GetLastTradePriceResponse>   
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>
</pre>
```

Using the RESTful facade is much easier:

<pre>
curl localhost:8080/quote/foo
34.5
</pre>
