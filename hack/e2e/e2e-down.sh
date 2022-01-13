#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

ROOT_DIR=$(dirname "${BASH_SOURCE[0]}")/../..
WORK_DIR=${ROOT_DIR}/testbin

cd ${WORK_DIR}
./kind delete cluster --name=devstream-e2e
