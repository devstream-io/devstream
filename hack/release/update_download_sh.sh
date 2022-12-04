#! /bin/bash -e

DOWNLOAD_SCRIPT_PATH="hack/install/download.sh"
STORAGE_BASE_URL=s3://download.devstream.io
DOWNLOAD_SCRIPT_S3_URL=${STORAGE_BASE_URL}/${DOWNLOAD_SCRIPT_PATH}

# upload download.sh to aws s3
echo "Uploading ${DOWNLOAD_SCRIPT_PATH} to ${DOWNLOAD_SCRIPT_S3_URL} ..."
aws s3 cp ${DOWNLOAD_SCRIPT_PATH} ${DOWNLOAD_SCRIPT_S3_URL} --acl public-read
