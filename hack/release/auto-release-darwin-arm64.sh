#! /bin/bash -e
# usage: `sh auto-release-darwin-arm64.sh -t v0.6.0`
set -o nounset

tag="invalid"

while getopts "t:" opt; do
  case $opt in
  t)
    tag=$OPTARG
    ;;

  ?)
    echo "Options not used"
    exit 1
    ;;
  esac
done

if [ "${tag}" == "invalid" ]; then
  echo "Maybe you forgot to use -t flag. E.g. sh auto-release-darwin-arm64.sh -t v0.6.0"
  exit 1
fi
echo "tag: ${tag}"

user=devstream-io
repo=devstream
github_token=$GITHUB_TOKEN
plugin_dir=.devstream

GOOS=$(go env GOOS)
GOARCH=$(go env GOARCH)

if [ ! $tag ] || [ ! $user ] || [ ! $repo ] || [ ! $github_token ] || [ ! $plugin_dir ]; then
  echo "The following variables cannot be empty!"
  echo "tag="$tag
  echo "user="$user
  echo "repo="$repo
  if [ ! $github_token ]; then
    echo "github_token="$github_token
  else
    echo "github_token=***"
  fi
  echo "plugin_dir="$plugin_dir
  exit
fi

# call build core and plugins
cd ../..
make build -j8

# install github-release for uploading
go install github.com/github-release/github-release@latest

# upload dtm
echo 'Uploading 'dtm-${GOOS}-${GOARCH}' ...'
github-release upload --security-token $github_token --user $user --repo $repo --tag $tag --file dtm --name dtm-${GOOS}-${GOARCH}
echo dtm-${GOOS}-${GOARCH}' uploaded.'

# upload plugins and .md5 files
# In order to upload plug-ins to s3, you need to download awscli. After downloading awscli, you need to configure aws credentials.
pip3 install awscli
aws s3 cp $plugin_dir s3://download.devstream.io/${tag} --recursive --acl public-read
