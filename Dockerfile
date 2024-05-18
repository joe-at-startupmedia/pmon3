FROM golang:1.21-alpine

RUN apk --update add build-base && \
  apk add --no-cache git && \
  mkdir -p /usr/src/pmon3 && \
  cd /usr/src/pmon3 && \
  git clone https://github.com/joe-at-startupmedia/pmon3.git . && \
  mkdir /usr/src/pmon3/data && \
  mkdir /usr/src/pmon3/logs && \
  make test && \
  make test_cgo
