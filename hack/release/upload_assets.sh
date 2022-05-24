#! /bin/bash -e

github_token=$1
tag=$2
GOOS=$3
GOARCH=$4
user=devstream-io
repo=devstream
plugin_dir=.devstream

# upload dtm
echo 'Uploading 'dtm-${GOOS}-${GOARCH}' ...'
github-release upload --security-token $github_token --user $user --repo $repo --tag $tag --file dtm --name dtm-${GOOS}-${GOARCH}
echo dtm-${GOOS}-${GOARCH}' uploaded.'

