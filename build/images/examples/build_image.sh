#!/bin/bash
#
# Build Docker images for Kubebench examples.
# This is intended to be invoked as a step in Argo to build the docker images.
#
# The script checks all subdirectory for existence of build_image.sh scripts.
# If so, call it for each *.Dockerfile in that subdirectory.
#
# build_image.sh ${SRC_DIR} ${DOCKERFILE_ROOT} ${IMAGE_PREFIX} ${VERSION}
set -ex

SRC_DIR=$(realpath $1)
DOCKERFILE_ROOT=$(realpath $2)
IMAGE_PREFIX=$3
VERSION=$4
if [ -z ${VERSION} ]; then
  VERSION=$(git describe --tags --always --dirty)
fi

for dir in $(find ${DOCKERFILE_ROOT} -mindepth 1 -maxdepth 1 -type d); do
  if [ -e ${dir}/build_image.sh ]; then
    echo "Build images from ${dir}"
    for dockerfile in ${dir}/*.Dockerfile; do
      filename=$(basename $dockerfile)
      image=${IMAGE_PREFIX}-$(basename $dir)-${filename%.Dockerfile}
      /bin/bash ${dir}/build_image.sh ${SRC_DIR} ${dockerfile} ${image} ${VERSION}
    done
  fi
done

echo "All images built successfully."
