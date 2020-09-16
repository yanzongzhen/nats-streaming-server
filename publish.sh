!/bin/bash

appName="nats-streaming-server"

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build &&
scp $appName H1:/data/server && echo "H1 OK" &&
scp $appName H2:/data/server && echo "H2 OK" &&
scp $appName H3:/data/server && echo "H3 OK" &&
rm -rf $appName && echo "clean OK"
