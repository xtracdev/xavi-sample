#!/bin/bash

curl -i -X POST -H 'Content-Type: application/json' -d@imposter.json http://127.0.0.1:2525/imposters

export XAVI_KVSTORE_URL=file:///`pwd`/config
./xavisample add-server -address localhost -port 4545 -name quotesvr1 -health-check-interval 5000 -health-check-timeout 2000 -health-check custom-http -ping-uri /
./xavisample add-server -address localhost -port 4646 -name quotesvr2 -health-check-interval 5000 -health-check-timeout 2000 -health-check custom-http -ping-uri /
./xavisample add-server -address localhost -port 4747 -name quotesvr3 -health-check-interval 5000 -health-check-timeout 2000 -health-check custom-http -ping-uri /
./xavisample add-backend -name quote-backend -servers quotesvr1,quotesvr2,quotesvr3 -load-balancer-policy prefer-local
./xavisample add-route -name quote-route -backends quote-backend -base-uri /quote/ -plugins Quote,SessionId,Timing,Recovery
./xavisample add-listener -name quote-listener -routes quote-route

./xavisample listen -ln quote-listener -address 0.0.0.0:8080
