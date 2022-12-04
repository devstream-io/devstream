#! /bin/bash -e

version=$1

LATEST_VERSION_FILE="latest_version"
STORAGE_BASE_URL=s3://download.devstream.io

# create latest_version file
echo "Saving latest version(${version}) to ${LATEST_VERSION_FILE} ..."
echo $version > ${LATEST_VERSION_FILE}

# create or update latest_version.txt on s3
aws s3 cp ${LATEST_VERSION_FILE} ${STORAGE_BASE_URL}/${LATEST_VERSION_FILE} --acl public-read
echo "${LATEST_VERSION_FILE} uploaded."
