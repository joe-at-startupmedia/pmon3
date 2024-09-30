FROM golang:1.22-alpine

ARG GIT_BRANCH_ARG=master
ARG MAKE_TARGET_ARG=test

ENV GIT_BRANCH=${GIT_BRANCH_ARG}
ENV MAKE_TARGET=${MAKE_TARGET_ARG}



RUN apk --update add build-base && \
  apk add --no-cache git && \
  mkdir -p /usr/src/pmon3 && \
  cd /usr/src/pmon3 && \
  git clone --single-branch --branch "$GIT_BRANCH" https://github.com/joe-at-startupmedia/pmon3.git . && \
  mkdir /usr/src/pmon3/data && \
  mkdir /usr/src/pmon3/logs && \
  make "$MAKE_TARGET"