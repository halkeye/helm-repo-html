#!/bin/bash

# Combination of the Glide and Helm scripts (Borrowed from https://github.com/bacongobbler/helm-whatup/blob/master/install-binary.sh)

PROJECT_NAME="helm-repo-html"
PROJECT_GH="halkeye/$PROJECT_NAME"

[ -z "$HELM_HOME" ] && HELM_HOME=$(helm env | grep 'HELM_DATA_HOME' | cut -d '=' -f    2 | tr -d '"')
: ${HELM_PLUGIN_PATH:="$HELM_HOME/plugins/helm-repo-html"}

if [[ $SKIP_BIN_INSTALL == "1" ]]; then
  echo "Skipping binary install"
  exit
fi

# initArch discovers the architecture for this system.
initArch() {
  ARCH=$(uname -m)
  case $ARCH in
    armv5*) ARCH="armv5";;
    armv6*) ARCH="armv6";;
    armv7*) ARCH="armv7";;
    aarch64) ARCH="arm64";;
    x86) ARCH="i386";;
    x86_64) ARCH="x86_64";;
    i686) ARCH="i386";;
    i386) ARCH="i386";;
  esac
}

# initOS discovers the operating system for this system.
initOS() {
  OS=$(echo `uname`|tr '[:upper:]' '[:lower:]')

  case "$OS" in
    # Minimalist GNU for Windows
    mingw*) OS='windows';;
  esac
}

# verifySupported checks that the os/arch combination is supported for
# binary builds.
verifySupported() {
  if ! type "curl" > /dev/null && ! type "wget" > /dev/null; then
    echo "Either curl or wget is required"
    exit 1
  fi
}

# getDownloadURL checks the latest available version.
getDownloadURL() {
  # Use the GitHub API to find the latest version for this project.
  local latest_url="https://api.github.com/repos/$PROJECT_GH/releases/latest"
  if type "curl" > /dev/null; then
    DOWNLOAD_URL=$(curl -s $latest_url | grep $OS | grep $ARCH | awk '/\"browser_download_url\":/{gsub( /[,\"]/,"", $2); print $2}')
  elif type "wget" > /dev/null; then
    DOWNLOAD_URL=$(wget -q -O - $latest_url | grep $OS | grep $ARCH | awk '/\"browser_download_url\":/{gsub( /[,\"]/,"", $2); print $2}')
  fi
}

# downloadFile downloads the latest binary package and also the checksum
# for that binary.
downloadFile() {
  PLUGIN_TMP_FILE="/tmp/${PROJECT_NAME}.tgz"
  echo "Downloading $DOWNLOAD_URL"
  if type "curl" > /dev/null; then
    curl -s -L "$DOWNLOAD_URL" -o "$PLUGIN_TMP_FILE"
  elif type "wget" > /dev/null; then
    wget -q -O "$PLUGIN_TMP_FILE" "$DOWNLOAD_URL"
  fi
}

# installFile unpacks and installs plugin
installFile() {
  HELM_TMP="/tmp/$PROJECT_NAME"
  mkdir -p "$HELM_TMP"
  tar xf "$PLUGIN_TMP_FILE" -C "$HELM_TMP"
  echo "Preparing to install into ${HELM_PLUGIN_PATH}"
  cp -R "$HELM_TMP/bin" "$HELM_PLUGIN_PATH/"
}

# fail_trap is executed if an error occurs.
fail_trap() {
  result=$?
  if [ "$result" != "0" ]; then
    echo "Failed to install $PROJECT_NAME"
    echo "\tFor support, go to https://github.com/${PROJECT_GH}"
  fi
  exit $result
}

# testVersion tests the installed client to make sure it is working.
testVersion() {
  set +e
  echo "$PROJECT_NAME installed into $HELM_PLUGIN_PATH/$PROJECT_NAME"
  $HELM_PLUGIN_PATH/bin/$PROJECT_NAME -h
  set -e
}

# Execution

#Stop execution on any error
trap "fail_trap" EXIT
set -e
initArch
initOS
verifySupported
getDownloadURL
downloadFile
installFile
testVersion
