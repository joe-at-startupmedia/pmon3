#!/usr/bin/env bash

set -ex

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

systemd_install() {
  WHOAMI=$(whoami)
  sudo cp -R bin/pmon* /usr/local/bin/
  sudo cp "rpm/pmond.service" /usr/lib/systemd/system/
  sudo cp "rpm/pmond.logrotate" /etc/logrotate.d/pmond
  sudo mkdir -p /var/log/pmond/ /etc/pmon3/config/ /etc/pmon3/data/
  # prevent configuration overwrite from previous installation
  if [ ! -f /etc/pmon3/config/config.yml ]; then
    sudo cp "config.yml" /etc/pmon3/config/
  fi
  sudo systemctl enable pmond
  sudo systemctl start pmond
  sleep 2
  sudo sh -c "bin/pmon3 completion bash > /etc/profile.d/pmon3.sh"
  sudo chown -R "root:${WHOAMI}" /var/log/pmond
  sudo chmod 0755 /var/log/pmond/ || true
  ./bin/pmon3 ls
  ./bin/pmon3 --help
}


echo "Installing $PROJECT_NAME from release: $RELEASE"

#the extracted folder isnt prepended by the letter v
download_from_project "$PROJECT_URL/archive/refs/tags/v$RELEASE.tar.gz" "$RELEASE_ARCHIVE.tar.gz"


rm -rf "$RELEASE_ARCHIVE" && \
  tar -xvzf "$RELEASE_ARCHIVE.tar.gz" && \
  rm -f "$RELEASE_ARCHIVE.tar.gz" && \
  cd "$RELEASE_ARCHIVE" && \
  mkdir bin && \
  download_from_project "$PROJECT_URL/releases/download/v$RELEASE/pmon3" "bin/pmon3" && \
  download_from_project "$PROJECT_URL/releases/download/v$RELEASE/pmond" "bin/pmond" && \
  chmod +x bin/* && \
  systemd_install

clean_downloads
