FROM golang:1.23.2-bullseye

ARG GIT_BRANCH_ARG=master
ARG MAKE_TARGET_ARG=test
ARG TEST_REGEX_ARG=Test
ARG TEST_PACKAGES_ARG
ARG CODECOV_TOKEN_ARG

RUN  apt-get update && \
  apt-get install -y build-essential git curl bash jq && \
  curl -L -o /usr/local/bin/codecov https://github.com/codecov/codecov-cli/releases/download/v0.7.5/codecovcli_linux && \
  chmod +x /usr/local/bin/codecov && \
  ls /usr/local/bin/ && \
  cd /opt/ && \
  git clone --single-branch --branch ${GIT_BRANCH_ARG} https://github.com/joe-at-startupmedia/pmon3.git

ENV CODECOV_TOKEN=${CODECOV_TOKEN_ARG}
ENV TEST_REGEX=${TEST_REGEX_ARG}
ENV MAKE_TARGET=${MAKE_TARGET_ARG}
ENV TEST_PACKAGES=${TEST_PACKAGES_ARG}

ENTRYPOINT ["/bin/sh", "-c" , "cd /opt/pmon3 && make build && make install && CODECOV_TOKEN=${CODECOV_TOKEN} CODECOV_FLAG=${MAKE_TARGET} make ${MAKE_TARGET}"]