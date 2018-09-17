#!/bin/bash
#
# Build Docker images for Kubebench example tf-cnn.
#
# build_image.sh ${DOCKERFILE} ${IMAGE} ${TAG}
set -ex

DOCKERFILE=$(realpath $1)
IMAGE=$2
TAG=$3

DOCKERFILE_DIR=$(dirname $DOCKERFILE)
SRC_ROOT=${DOCKERFILE_DIR%/build/images/examples/tf-cnn}

echo "Authenticate gcloud account"
gcloud auth activate-service-account --key-file=${GOOGLE_APPLICATION_CREDENTIALS}

cd $SRC_ROOT
echo "Build image ${IMAGE}:${TAG}"
docker build -t ${IMAGE}:${TAG} -f ${DOCKERFILE} .
echo "Push image ${IMAGE}:${TAG}"
gcloud docker -- push "${IMAGE}:${TAG}"
echo "Image ${IMAGE}:${TAG} built successfully"
