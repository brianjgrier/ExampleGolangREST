#~/bin/sh

CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -tags netgo -ldflags '-w -extldflags "-static"' -o main .

docker build -t brianjgrier/simple_rest_sql:latest . --no-cache -f Dockerfile

docker push brianjgrier/simple_rest_sql:latest

