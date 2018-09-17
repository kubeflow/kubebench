#!/bin/bash
#
# Build Docker images for Kubebench examples.
# This is intended to be invoked as a step in Argo to build the docker images.
#
# The script checks all subdirectory for existence of build_image.sh scripts.
# If so, call it for each *.Dockerfile in that subdirectory.
#
# build_image.sh ${DOCKERFILE_PLACEHOLDER} ${IMAGE_PREFIX} ${TAG}
#
#  $DOCKERFILE_PLACEHOLDER:
#    Path to a placeholder (non-existent) file under same directory as this script.
#    This parameter is used to get parent directory info.
#  $IMAGE_PREFIX:
#    The name prefix for all images built by this script, the final image name is
#    in such format: $IMAGE_PREFIX-$SUBDIRECTORY_NAME-$FILENAME_WITHOUT_EXTENSION.
#  $TAG:
#    The image tag to be used for all images built by this script.
set -ex

DOCKERFILE_PLACEHOLDER=$(realpath $1)
IMAGE_PREFIX=$2
TAG=$3

BUILD_FILE_DIR=$(dirname $DOCKERFILE_PLACEHOLDER)

for dir in $(find ${BUILD_FILE_DIR} -mindepth 1 -maxdepth 1 -type d); do
  if [ -e ${dir}/build_image.sh ]; then
    echo "Build images from ${dir}"
    for dockerfile in ${dir}/*.Dockerfile; do
      filename=$(basename $dockerfile)
      image=${IMAGE_PREFIX}-$(basename $dir)-${filename%.Dockerfile}
      /bin/bash ${dir}/build_image.sh ${dockerfile} ${image} ${TAG}
    done
  fi
done

echo "All images built successfully."
