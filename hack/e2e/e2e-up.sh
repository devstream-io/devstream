#!/usr/bin/env bash

# This script is mainly used to start a k8s 1-node cluster environment for e2e testing.
# Docker is the only dependency that needs to be prepared in advance.
# `kubectl` and `kind` will be obtained through `curl`, maybe this way has the best compatibility.

# All materials required for testing are downloaded to the testbin directory,
# which is initialized only when it is executed for the first time, and then cached locally.
# The reason for not judging whether there are related tools like kubectl locally
# is because these tools may not match the versions required and are not easy to manage.

set -o errexit
set -o nounset
set -o pipefail

ROOT_DIR=$(dirname "${BASH_SOURCE[0]}")/../..
WORK_DIR=${ROOT_DIR}/testbin
K8S_VERSION=${TESTENV_K8S_VERSION:-1.22.0}
KIND_VERSION=${TESTENV_KIND_VERSION:-0.11.1}

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
    HOST_ARCH="arm64"
  elif [ "$(arch)" == "arm64" ];then
    HOST_ARCH="amd64"
  else
    echo "Support amd64/arm64 CPU arch only"
    exit 1
  fi

  echo "Got OS type: ${HOST_OS} and CPU arch: ${HOST_ARCH}"
}

# docker
function check_deps() {
  echo "Requires having Docker installed"
  docker ps > /dev/null
  if [ $? -ne 0 ]; then
    exit 1
  fi
}

# kubectl & kind
function fetch_tools() {
  echo "All tools save in ${WORK_DIR} directory"
  if [ ! -f ${WORK_DIR} ]; then
    mkdir -p ${WORK_DIR}
  fi
  cd ${WORK_DIR}

  if [ ! -f kind ]; then
    echo "kind doesn't exist, download it now"
    kind_uri="https://kind.sigs.k8s.io/dl/v${KIND_VERSION}/kind-${HOST_OS}-${HOST_ARCH}"
    echo "kind uri: ${kind_uri}"
    curl -Lo ./kind "${kind_uri}"
  fi
  echo "kind already downloaded"
  chmod +x ./kind

  if [ ! -f kubectl ]; then
    echo "kubectl doesn't exist, download it now"
    kubectl_uri="https://dl.k8s.io/release/v${K8S_VERSION}/bin/${HOST_OS}/${HOST_ARCH}/kubectl"
    echo "kubectl uri: ${kubectl_uri}"
    curl -Lo ./kubectl "${kubectl_uri}"
  fi
  echo "kubectl already downloaded"
  chmod +x ./kubectl

  KUBECTL="${WORK_DIR}/kubectl"
  KIND="${WORK_DIR}/kind"

  cd -
}

# start up a kind cluster
function setup_env() {
  docker pull "kindest/node:v${K8S_VERSION}"
  ${KIND} create cluster --image="kindest/node:v${K8S_VERSION}" --config="${ROOT_DIR}/hack/e2e/kind.yaml" --name=devstream-e2e
}

init
check_deps
fetch_tools
setup_env
