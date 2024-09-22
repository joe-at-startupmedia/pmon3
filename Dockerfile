FROM golang:1.22-alpine

ENV GIT_BRANCH=${GIT_BRANCH}

RUN apk --update add build-base && \
  apk add --no-cache git && \
  mkdir -p /usr/src/pmon3 && \
  cd /usr/src/pmon3 && \
  git clone -b "$GIT_BRANCH" https://github.com/joe-at-startupmedia/pmon3.git . && \
  mkdir /usr/src/pmon3/data && \
  mkdir /usr/src/pmon3/logs && \
  make test && \
  make test_cgo
