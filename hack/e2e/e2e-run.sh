#!/usr/bin/env bash

# set -o errexit
set -o nounset
set -o pipefail

ROOT_DIR=$(dirname "${BASH_SOURCE[0]}")/../..
SCRIPT_DIR=${ROOT_DIR}/hack/e2e
CONFIG_DIR=${ROOT_DIR}/test/e2e/yaml

READY=0
TIME=0
TIMEOUT=120

# setup k8s cluster and setup config yaml files
function init() {
    echo "[dtm e2e test script] Create k8s cluster by kind!"
    bash ${SCRIPT_DIR}/e2e-up.sh

    cp ${CONFIG_DIR}/* ${ROOT_DIR}
}

function dtm_apply() {
    cd ${ROOT_DIR}

    echo "[dtm e2e test script] Start dtm apply!"
    ./dtm apply -f e2e-config-local.yaml -y
}

# timeout set to 2 min
function check_pod_status() {
    echo "[dtm e2e test script] Start check pod status!"
    while [ "$(kubectl get pods -l app=dtm-e2e-go -o 'jsonpath={..status.conditions[?(@.type=="Ready")].status}')" != "True" ]; do
        echo "pod not ready yet..."
        sleep 5
        TIME=$((TIME + 5))
        if [ ${TIME} -ge ${TIMEOUT} ]; then
            return
        fi
    done
    READY=1
}

function check_ready() {
    if [ ${READY} -eq 1 ]; then
        echo "[dtm e2e test script] Pod is ready!"
        echo "[dtm e2e test script] E2E test success!"
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
    rm -rf ${ROOT_DIR}/e2e*.yaml
}

case "$1" in
run)
    init
    dtm_apply
    check_pod_status
    check_ready
    # clean
    ;;
clean)
    clean
    ;;
*)
    echo "Usage:$0 somethine.Please see description!"
    ;;
esac
