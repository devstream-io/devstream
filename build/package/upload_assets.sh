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

# upload each plugin
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
