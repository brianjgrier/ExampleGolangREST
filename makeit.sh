#~/bin/sh

CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -tags netgo -ldflags '-w -extldflags "-static"' -o main .

docker build -t sds-server1:32000/simple_rest_sql:registry . --no-cache -f Dockerfile

docker push sds-server1:32000/simple_rest_sql