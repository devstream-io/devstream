#! /bin/bash -e
# usage: `sh auto-release-darwin-arm64.sh -t v0.6.0`
set -o nounset
# exit when any command fails
set -e

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
plugin_dir=~/.devstream/plugins

GOOS=$(go env GOOS)
GOARCH=$(go env GOARCH)
DTM_CORE_BINARY=dtm-${GOOS}-${GOARCH}
STORAGE_BASE_URL=s3://download.devstream.io
STORAGE_URL_WITH_TAG=${STORAGE_BASE_URL}/${tag}

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
make clean
make build -j8

# install github-release for uploading
go install github.com/github-release/github-release@latest

# upload dtm
echo "Uploading ${DTM_CORE_BINARY} ..."
# upload dtm to github release
github-release upload --security-token $github_token --user $user --repo $repo --tag $tag --file dtm --name ${DTM_CORE_BINARY}
# upload dtm to aws s3
aws s3 cp dtm ${STORAGE_URL_WITH_TAG}/${DTM_CORE_BINARY} --acl public-read
echo "${DTM_CORE_BINARY} uploaded."

# upload plugins and .md5 files
# In order to upload plug-ins to s3, you need to download aws cli.
# After downloading aws cli, you need to configure aws credentials.
pip3 install awscli
aws s3 cp $plugin_dir $STORAGE_URL_WITH_TAG --recursive --acl public-read

# check if the number of plugins on s3 is correct
local_plugin_nums=$(../../dtm list plugins |wc -l)
((local_plugin_file_nums=local_plugin_nums*6))
s3_plugin_file_total_nums=$(aws s3 ls download.devstream.io/"$tag"/|awk '{print $NF}'|uniq|wc -l)
((s3_plugin_file_nums=s3_plugin_file_total_nums-3))
echo "s3_plugin_file_nums:" "$s3_plugin_file_nums"
echo "local_plugin_file_nums:" "$local_plugin_file_nums"
if [ "$local_plugin_file_nums" -ne "$s3_plugin_file_nums" ]
then
  echo "Attention,Maybe the plugin uploaded to s3 is not correct."
else
  echo "The plugin uploaded to s3 is correct."
fi
