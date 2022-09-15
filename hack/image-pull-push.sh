#!/usr/bin/env bash

CurrentDIR=$(cd "$(dirname "$0")" || exit;pwd)

func() {
    echo "\"image-pull-push.sh\" is used for pull images from dockerhub and push images to private image registry."
    echo ""
    echo "Usage:"
    echo "1. Download images from dockerhub and save them to a tarball:"
    echo "./image-pull-push.sh -f images-list.txt -d dtm-images -r harbor.devstream.io -s"
    echo ""
    echo "2. Load images from tarball locally:"
    echo "./image-pull-push.sh -f images-list.txt -d dtm-images -l"
    echo ""
    echo "3. Upload images to private image registry:"
    echo "./image-pull-push.sh -f images-list.txt -d dtm-images -r harbor.devstream.io -u"
    echo ""
    echo "4. Load and upload at the same time:"
    echo "./image-pull-push.sh -f images-list.txt -d dtm-images -r harbor.devstream.io -l -u"
    echo ""
    echo "5. Short command with default config: ./image-pull-push.sh -f images-list.txt -l -u"
    echo ""
    echo "  $0 [-f IMAGES-LIST-FILE] [-d IMAGES-DIR] [-r PRIVATE-REGISTRY] [-s | -l | -u]"
    echo ""
    echo "Description:"
    echo "  -f IMAGES-LIST-FILE    : text file with image list."
    echo "  -r PRIVATE-REGISTRY    : target private registry addr. (default: harbor.devstream.io)"
    echo "  -d IMAGES-DIR          : a directory where images.tar.gz is saved. (default: dtm-images)"
    echo "  -s                     : save images"
    echo "  -l                     : load images"
    echo "  -u                     : upload images"
    echo ""
    echo "Notice: pigz should be installed first. (yum install pigz)"
    exit
}

save='false'
load='false'
upload='false'

while getopts 'f:r:d:sluh' OPT; do
    case $OPT in
        f) ImagesListFile="$OPTARG";;
        r) RegistryUrl="$OPTARG";;
        d) ImagesDir="$OPTARG";;
        s) save='true';;
        l) load='true';;
        u) upload='true';;
        h) func;;
        ?) func;;
        *) func;;
    esac
done

if [[ -z $RegistryUrl ]]; then
  RegistryUrl='harbor.devstream.io'
fi

if [[ -z $ImagesDir ]]; then
  ImagesDir='dtm-images'
fi

# save
if [[ $save == 'true' ]]; then
    if [ ! -d ${ImagesDir} ]; then
       mkdir -p ${ImagesDir}
    fi
    ImagesListLen=$(cat ${ImagesListFile} | grep -v ^$ | awk '{print}' | wc -l)
    name=""
    images=""
    index=0
    for image in $(<${ImagesListFile}); do
        if [[ ${image} =~ ^\#\#.* ]]; then
           if [[ -n ${images} ]]; then
              echo ""
              echo "Save images: "${name}" to "${ImagesDir}"/"${name}".tar.gz  <<<"
              docker save ${images} | pigz > ${ImagesDir}"/"${name}.tar.gz
              echo ""
           fi
           images=""
           name=$(echo "${image}" | sed 's/#//g' | sed -e 's/[[:space:]]//g')
           ((index++))
           continue
        fi

        if [[ ${image} =~ ^----.* ]]; then
          ((index++))
           continue
        fi

        docker pull "${image}"
        docker tag "${image}" "${RegistryUrl}/${image}"
        images=${images}" "${RegistryUrl}/${image}
        if [[ ${index} -eq ${ImagesListLen}-1 ]]; then
           if [[ -n ${images} ]]; then
              echo "Save images: "${name}" to "${ImagesDir}"/"${name}".tar.gz  <<<"
              echo "Images: ${images}"
              docker save ${images} | pigz > ${ImagesDir}"/"${name}.tar.gz
           fi
        fi
        ((index++))
    done
fi

# load
if [[ $load == 'true' ]]; then
  for image in $(<${ImagesListFile}); do
    if [[ ${image} =~ ^\#\#.* ]]; then
      docker load -i ${ImagesDir}/$(echo "${image}" | sed 's/#//g' | sed -e 's/[[:space:]]//g').tar.gz
    fi
  done
fi

# push
if [[ $upload == 'true' ]]; then
  for image in $(<${ImagesListFile}); do
    if [[ ${image} =~ ^\#\#.* ]] || [[ ${image} =~ ^----.* ]]; then
      continue
    fi
    docker push "${RegistryUrl}/${image}"
  done
fi
