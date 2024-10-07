FROM golang:1.22-alpine

ARG GIT_BRANCH_ARG=master
ARG MAKE_TARGET_ARG=test
ARG TEST_REGEX_ARG=Test
ARG CODECOV_TOKEN_ARG

RUN apk --update add build-base && \
  apk add --no-cache git curl && \
  curl -o  /usr/local/bin/codecov https://cli.codecov.io/latest/alpine/codecov && \
  chmod +x /usr/local/bin/codecov && \
  ls -al /usr/local/bin/codecov && \
  mkdir -p /usr/src/pmon3 && \
  cd /usr/src/pmon3 && \
  git clone --single-branch --branch ${GIT_BRANCH_ARG} https://github.com/joe-at-startupmedia/pmon3.git . && \
  mkdir /usr/src/pmon3/data && \
  mkdir /usr/src/pmon3/logs  \

ENV TEST_REGEX=${TEST_REGEX_ARG}

ENTRYPOINT ["/bin/sh", "-c" , "make ${MAKE_TARGET_ARG} && /usr/local/bin/codecov upload-process -t ${CODECOV_TOKEN_ARG} -F ${MAKE_TARGET_ARG}"]