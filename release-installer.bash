#!/usr/bin/env bash

set -x

RELEASE=$1
PROJECT_NAME="pmon3"
PROJECT_URL="https://github.com/joe-at-startupmedia/pmon3"
RELEASE_ARCHIVE="$PROJECT_NAME-$RELEASE"

if [ "$RELEASE" == "" ]; then
    echo "Please enter release version"
    exit 1
fi

clean_downloads() {
  rm -f "$RELEASE_ARCHIVE.tar.gz"
}

download_from_project() {
  SRC=$1
  DST=$2
  wget "$SRC" -O "$DST"
  wgetreturn=$?
  if [[ $wgetreturn -ne 0 ]]; then
    echo "Could not wget: $SRC"
    clean_downloads
    exit 1
  fi
}


echo "Installing $PROJECT_NAME from release: $RELEASE"

#the extracted folder isnt prepended by the letter v
download_from_project "$PROJECT_URL/archive/refs/tags/v$RELEASE.tar.gz" "$RELEASE_ARCHIVE.tar.gz"

tar -xvzf "$RELEASE_ARCHIVE.tar.gz" &&
  rm -f "$RELEASE_ARCHIVE.tar.gz" && \
  cd "$RELEASE_ARCHIVE" && \
  mkdir bin && \
  download_from_project "$PROJECT_URL/releases/download/v$RELEASE/pmon3" "bin/pmon3" && \
  download_from_project "$PROJECT_URL/releases/download/v$RELEASE/pmond" "bin/pmond" && \
  chmod +x bin/* && \
  make systemd_install

clean_downloads
