function init() {
  if [ "$(uname)" == "Darwin" ];then
    HOST_OS="darwin"
  elif [ "$(uname)" == "Linux" ];then
    HOST_OS="linux"
  else
    echo "Support Darwin/Linux OS only"
    exit 1
  fi

  if [ "$(arch)" == "amd64" ];then
    HOST_ARCH="amd64"
  elif [ "$(arch)" == "arm64" ];then
    HOST_ARCH="arm64"
  else
    echo "Support amd64/arm64 CPU arch only"
    exit 1
  fi

  echo "Got OS type: ${HOST_OS} and CPU arch: ${HOST_ARCH}"
}

function getLatestReleaseVersion() {
  latestVersion=$(curl -s https://api.github.com/repos/devstream-io/devstream/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
  if [ -z "$latestVersion" ]; then
    echo "Failed to get latest release version"
    exit 1
  fi
  echo "Latest dtm release version: ${latestVersion}\n"
}

function downloadDtm() {
  # 1. download the release and rename it to "dtm"
  # 2. count the download count of the release
  fullReleaseUrl="https://devstream.gateway.scarf.sh/releases/$latestVersion/dtm-$HOST_OS-$HOST_ARCH"
  echo "Downloading dtm from: $fullReleaseUrl"
  # use -L to follow redirects
  curl -L -o dtm $fullReleaseUrl
  echo "dtm downloaded completed\n"

  # grant execution rights
  chmod +x dtm
}

function downloadQuickStartConfig() {
  # config file is small so we use -s to ignore the output
  curl -s -o quickstart.yaml https://raw.githubusercontent.com/devstream-io/devstream/main/examples/quickstart.yaml
  echo "quickstart.yaml downloaded completed"
}

function showDtmHelp() {
  echo ""
  # show dtm help and double check the download is success
  ./dtm help
}

init
getLatestReleaseVersion
downloadDtm
downloadQuickStartConfig
showDtmHelp

