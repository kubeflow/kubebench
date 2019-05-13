#!/bin/bash
#
# Build Docker images for Kubebench controller.
# This is intended to be invoked as a step in Argo to build the docker image.
#
# build_image.sh ${SRC_DIR} ${DOCKERFILE} ${IMAGE} ${VERSION}
set -ex

SRC_DIR=$(realpath $1)
DOCKERFILE=$(realpath $2)
IMAGE=$3
VERSION=$4
if [ -z ${VERSION} ]; then
  VERSION=$(git describe --tags --always --dirty)
fi
TAG=${REGISTRY}/${REPO_NAME}/${IMAGE}:${VERSION}

echo "Setup build directory"
export GOPATH=`mktemp -d -p $(dirname $SRC_DIR)`
export PATH=${GOPATH}/bin:/usr/local/go/bin:${PATH}
mkdir -p ${GOPATH}/src/github.com/kubeflow/kubebench
BUILD_DIR=${GOPATH}/src/github.com/kubeflow/kubebench

# need to download all packages into gopath (this options allows go build to download all data)
export GO111MODULE=on

echo "Copy source and Dockerfile to build directory"
cp -r ${SRC_DIR}/controller ${BUILD_DIR}/controller
cp ${DOCKERFILE} ${BUILD_DIR}/Dockerfile
cp  ${SRC_DIR}/go.mod ${BUILD_DIR}/go.mod
cp  ${SRC_DIR}/go.sum ${BUILD_DIR}/go.sum
echo "Change working directory to ${BUILD_DIR}"
cd ${BUILD_DIR}

echo "Build go binaries"
GOOS=linux CGO_ENABLED=0 go build github.com/kubeflow/kubebench/controller/cmd/configurator
GOOS=linux CGO_ENABLED=0 go build github.com/kubeflow/kubebench/controller/cmd/reporter

echo "Authenticate gcloud account"
gcloud auth activate-service-account --key-file=${GOOGLE_APPLICATION_CREDENTIALS}
echo "Build image ${TAG}"
gcloud builds submit --tag=${TAG} --project=${PROJECT} .

echo "Clean up go build directory"
cd
rm -rf ${GOPATH}

echo "Image ${TAG} built successfully"
