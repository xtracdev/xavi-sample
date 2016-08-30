#!/bin/bash

#Run server.js instead of mb to simulate timeout after headers written

export XAVI_KVSTORE_URL=file:///`pwd`/config
./xavisample add-server -address localhost -port 4545 -name quotesvr1 -health-check-interval 2000 -health-check-timeout 500 -health-check http-get -ping-uri /
./xavisample add-backend -name quote-backend -servers quotesvr1
./xavisample add-route -name quote-route -backends quote-backend -base-uri /quote/ -plugins Quote,SessionId,Timing,Recovery
./xavisample add-listener -name quote-listener -routes quote-route

./xavisample listen -ln quote-listener -address 0.0.0.0:8080
