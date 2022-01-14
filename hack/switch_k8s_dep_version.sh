#!/usr/bin/env bash

# https://github.com/kubernetes/kubernetes/issues/79384#issuecomment-521493597
# sh hack/switch_k8s_dep_version.sh v1.22.2

set -euo pipefail

VERSION=${1#"v"}
if [ -z "$VERSION" ]; then
  echo "Must specify version!"
  exit 1
fi
echo "Kubernetes version: v${VERSION}"

MODS=($(
  curl -sS https://raw.githubusercontent.com/kubernetes/kubernetes/v${VERSION}/go.mod |
  sed -n 's|.*k8s.io/\(.*\) => ./staging/src/k8s.io/.*|k8s.io/\1|p'
))
echo "Downloaded kubernetes/go.mod"

for MOD in "${MODS[@]}"; do
  V=$(
    go mod download -json "${MOD}@kubernetes-${VERSION}" |
    sed -n 's|.*"Version": "\(.*\)".*|\1|p'
  )
  go mod edit "-replace=${MOD}=${MOD}@${V}"
done

go get "k8s.io/kubernetes@v${VERSION}"
