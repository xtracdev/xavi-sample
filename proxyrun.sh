#!/bin/bash

export XAVI_KVSTORE_URL=file:///`pwd`/config
./xavisample add-server -address localhost -port 8080 -name quotesvr
./xavisample add-backend -name proxy-backend -servers quotesvr
./xavisample add-route -name proxy-route -backends proxy-backend -base-uri / -plugins Timing,Recovery
./xavisample add-listener -name proxy-listener -routes proxy-route

./xavisample listen -ln proxy-listener -address 0.0.0.0:9090
