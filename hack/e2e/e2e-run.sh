#!/usr/bin/env bash

set -o pipefail

ROOT_DIR=$(dirname "${BASH_SOURCE[0]}")/../..
SCRIPT_DIR=${ROOT_DIR}/hack/e2e
CONFIG_DIR=${ROOT_DIR}/test/e2e/yaml
CONFIG_FILENAME=e2e-test-local.yaml
# CONFIG_FILENAME=e2e-config.yaml
# TOOLS_FILENAME=e2e-tools.yaml
# VARIABLES_FILENAME=e2e-variables.yaml

function check_variables() {

    if [ -z ${GITHUB_USER} ]; then
        echo "You have to set environment variable 'GITHUB_USER' first!"
        usage
        exit 1
    fi

    if [ -z ${GITHUB_TOKEN} ]; then
        echo "You have to set environment variable 'GITHUB_TOKEN' first!"
        usage
        exit 1
    fi

    if [ -z ${DOCKERHUB_USERNAME} ]; then
        echo "You have to set environment variable 'DOCKERHUB_USERNAME' first!"
        usage
        exit 1
    fi

    if [ -z ${DOCKERHUB_TOKEN} ]; then
        echo "You have to set environment variable 'DOCKERHUB_TOKEN' first!"
        usage
        exit 1
    fi
}

# setup k8s cluster and setup config yaml files
function init() {
    set -u
    echo "[dtm e2e test script] Create k8s cluster by kind!"
    bash ${SCRIPT_DIR}/e2e-down.sh
    bash ${SCRIPT_DIR}/e2e-up.sh

    gen_config
}

# generate temporary config files
function gen_config() {
    set -u

    # modify user's github name in e2e-test template config file and generate temporary config file in devstream root path
    sed -e "s/GITHUBUSERNAME/${GITHUB_USER}/g" ${CONFIG_DIR}/${CONFIG_FILENAME} >${ROOT_DIR}/${CONFIG_FILENAME}

}

# dtm e2e test
function dtm_test() {
    set -u
    cd ${ROOT_DIR}

    echo "[dtm e2e test script] Start dtm e2e test locally!"
    ./dtm apply -f ${CONFIG_FILENAME} -y
    check_status

    ./dtm verify -f ${CONFIG_FILENAME}
    ./dtm delete -f ${CONFIG_FILENAME} -y
    cd -
}

function check_status() {
    set -u
    pod_ready=1
    time=0
    timeout=120
    echo "[dtm e2e test script] Start check pod status!"
    while [ "$(kubectl get pods -l app=dtm-e2e-go -o 'jsonpath={..status.conditions[?(@.type=="Ready")].status}')" != "True" ]; do
        echo "pod not ready yet..."
        sleep 5
        time=$((time + 5))
        if [ ${time} -ge ${timeout} ]; then
            ${pod_ready}=0
            break
        fi
    done

    if [ ${pod_ready} -eq 1 ]; then
        echo "[dtm e2e test script] Pod is ready!"
    else
        echo "[dtm e2e test script] Pod is not ready!"
        echo "[dtm e2e test script] E2E test failed!"
        clean
        exit 1
    fi
}

function clean() {
    echo "[dtm e2e test script] Start to clean test environment and configuration files!"
    echo "[dtm e2e test script] Remove k8s cluster!"
    bash ${SCRIPT_DIR}/e2e-down.sh

    echo "[dtm e2e test script] Remove yaml files!"
    rm -rf ${ROOT_DIR}/${CONFIG_FILENAME}
    echo "[dtm e2e test script] E2E test success!"
}

function usage() {
    echo "Usage: bash $0.
  Before start e2e test locally, you have to set some environment variables, including:
  - 'GITHUB_USER'
  - 'GITHUB_TOKEN'
  - 'DOCKERHUB_USERNAME'
  - 'DOCKERHUB_TOKEN'."
}

check_variables
init
dtm_test
clean
