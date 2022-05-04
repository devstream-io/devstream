#! /bin/bash -e
set -o nounset

tag=v0.4.0

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

# In order to upload plug-ins to s3, you need to download awscli. After downloading awscli, you need to configure aws credentials.
pip3 install awscli -i https://pypi.tuna.tsinghua.edu.cn/simple/
aws s3 cp $plugin_dir  s3://download.devstream.io/${tag} --recursive --acl public-read

# upload dtm
echo 'Uploading 'dtm-${GOOS}-${GOARCH}' ...'
github-release upload --security-token $github_token --user $user --repo $repo --tag $tag --file dtm --name dtm-${GOOS}-${GOARCH}
echo dtm-${GOOS}-${GOARCH}' uploaded.'

# upload plugins and .md5 files
upload $plugin_dir
