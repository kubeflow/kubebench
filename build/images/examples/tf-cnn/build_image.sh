#!/bin/bash
#
# Build Docker images for Kubebench example tf-cnn.
#
# build_image.sh ${SRC_DIR} ${DOCKERFILE} ${IMAGE} ${VERSION}
set -ex

SRC_DIR=$(realpath $1)
DOCKERFILE=$(realpath $2)
IMAGE=$3
VERSION=$4
TAG=${REGISTRY}/${REPO_NAME}/${IMAGE}:${VERSION}

echo "Setup build directory"
export BUILD_DIR=`mktemp -d -p $(dirname $SRC_DIR)`

echo "Copy source and Dockerfile to build directory"
cp -r ${SRC_DIR}/examples ${BUILD_DIR}/examples
cp ${DOCKERFILE} ${BUILD_DIR}/Dockerfile

echo "Change working directory to ${BUILD_DIR}"
cd ${BUILD_DIR}

echo "Authenticate gcloud account"
gcloud auth activate-service-account --key-file=${GOOGLE_APPLICATION_CREDENTIALS}
echo "Build image ${TAG}"
gcloud builds submit --tag=${TAG} --project=${PROJECT} .

echo "Clean up build directory"
cd
rm -rf ${BUILD_DIR}

echo "Image ${TAG} built successfully"
