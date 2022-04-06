#! /bin/bash -e

tag=v0.3.2

user=devstream-io
repo=devstream
github_token=$GITHUB_TOKEN
plugin_dir=.devstream

GOOS=$(go env GOOS)
GOARCH=$(go env GOARCH)

# call build core and plugins
cd ../..
make build -j8

# install github-release for uploading
go install github.com/github-release/github-release@latest

# upload each plugin
function upload(){
for file in `ls $plugin_dir`
do
 if [ -d $plugin_dir"/"$file ]
 then
   read_dir $plugin_dir"/"$file
 else
   echo 'Uploading '$file' ...'
   github-release upload --security-token $github_token --user $user --repo $repo --tag $tag --file  $plugin_dir"/"$file --name $file
   echo $file' uploaded.'
 fi
done
}

# upload dtm
echo 'Uploading 'dtm-${GOOS}-${GOARCH}' ...'
github-release upload --security-token $github_token --user $user --repo $repo --tag $tag --file dtm --name dtm-${GOOS}-${GOARCH}
echo dtm-${GOOS}-${GOARCH}' uploaded.'

# upload plugins and .md5 files
upload $plugin_dir
