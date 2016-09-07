#!/bin/bash

curl -i -X POST -H 'Content-Type: application/json' -d@imposter-https.json http://127.0.0.1:2525/imposters

export XAVI_KVSTORE_URL=file:///`pwd`/config
../xavisample add-server -address `hostname` -port 4443 -name quotesvr1
../xavisample add-backend -name quote-backend -servers quotesvr1 -cacert-path ./cert.pem -tls-only=true
../xavisample add-route -name quote-route -backends quote-backend -base-uri /quote/ -plugins Quote,SessionId,Timing,Recovery
../xavisample add-listener -name quote-listener -routes quote-route
