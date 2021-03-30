# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# start from a the latest python image

FROM ubuntu:18.04
LABEL maintainer="Brian Grier <brinjgrier@gmail.com>"
LABEL version="0.1"

RUN apt-get update
RUN apt-get install -y ca-certificates
RUN update-ca-certificates

RUN useradd -ms /bin/bash appuser

WORKDIR /home/appuser

COPY main ./
RUN chown appuser:appuser main

USER appuser
ENV PATH="/home/appuser:${PATH}"

ENTRYPOINT [ "main" ]

EXPOSE 8080
