#!/usr/bin/env bash

set -ex

RELEASE=$1
PROJECT_NAME="pmon3"
AUTHOR_NAME="joe-at-startupmedia"
PROJECT_URL="https://github.com/$AUTHOR_NAME/$PROJECT_NAME"

if [ "$RELEASE" == "" ]; then
  RELEASE=$(wget -q -O - "https://api.github.com/repos/$AUTHOR_NAME/$PROJECT_NAME/tags" | jq -r '.[0].name')
fi

RELEASE_ARCHIVE="$PROJECT_NAME-$RELEASE"

clean_downloads() {
  rm -f "$RELEASE_ARCHIVE.tar.gz"
  rm -rf "${RELEASE_ARCHIVE//v}"
}

download_from_project() {
  SRC=$1
  DST=$2
  wget "$SRC" -O "$DST"
  wgetreturn=$?
  if [[ $wgetreturn -ne 0 ]]; then
    echo "Could not wget: $SRC"
    exit 1
  fi
}

systemd_install() {
  sudo systemctl stop pmond
  sudo cp -R bin/pmon* /usr/local/bin/
  sudo cp "rpm/pmond.service" /usr/lib/systemd/system/
  sudo cp "rpm/pmond.logrotate" /etc/logrotate.d/pmond
  # prevent configuration overwrite from previous installation
  if [ ! -f /etc/pmon3/config/config.yml ]; then
    sudo mkdir -p /etc/pmon3/config/
    sudo cp "config.yml" /etc/pmon3/config/
  fi
  sudo systemctl enable pmond
  sudo systemctl start pmond
  sleep 2
  sudo sh -c "bin/pmon3 completion bash > /etc/profile.d/pmon3.sh"
  ./bin/pmon3 ls
  ./bin/pmon3 --help
}


echo "Installing $PROJECT_NAME from release: $RELEASE"

#the extracted folder isn't prepended by the letter v
download_from_project "$PROJECT_URL/archive/refs/tags/$RELEASE.tar.gz" "$RELEASE_ARCHIVE.tar.gz"


rm -rf "$RELEASE_ARCHIVE" && \
  tar -xvzf "$RELEASE_ARCHIVE.tar.gz" && \
  rm -f "$RELEASE_ARCHIVE.tar.gz" && \
  cd "${RELEASE_ARCHIVE//v}" && \
  mkdir bin | true && \
  download_from_project "$PROJECT_URL/releases/download/$RELEASE/pmon3" "bin/pmon3" && \
  download_from_project "$PROJECT_URL/releases/download/$RELEASE/pmond" "bin/pmond" && \
  chmod +x bin/* && \
  systemd_install

clean_downloads
