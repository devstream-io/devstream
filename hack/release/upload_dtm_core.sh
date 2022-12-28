#! /bin/bash -e

github_token=$1
tag=$2
GOOS=$3
GOARCH=$4

DTM_CORE_BINARY=dtm-${GOOS}-${GOARCH}
STORAGE_BASE_URL=s3://download.devstream.io
STORAGE_URL_WITH_TAG=${STORAGE_BASE_URL}/${tag}

user=devstream-io
repo=devstream
plugin_dir=~/.devstream/plugins

# upload dtm core
echo "Uploading ${DTM_CORE_BINARY} ..."
# upload dtm to github release
github-release upload --security-token $github_token --user $user --repo $repo --tag $tag --file dtm --name ${DTM_CORE_BINARY}
# upload dtm to aws s3
aws s3 cp dtm ${STORAGE_URL_WITH_TAG}/${DTM_CORE_BINARY} --acl public-read
echo "${DTM_CORE_BINARY} uploaded."

