#!/bin/bash
#
# Build Docker images for Kubebench controller.
# This is intended to be invoked as a step in Argo to build the docker image.
#
# build_image.sh ${DOCKERFILE} ${IMAGE} ${TAG}
set -ex

DOCKERFILE=$(realpath $1)
IMAGE=$2
TAG=$3

DOCKERFILE_DIR=$(dirname $DOCKERFILE)
SRC_ROOT=${DOCKERFILE_DIR%/build/images/controller}

echo "Setup go build directory"
export GOPATH=`mktemp -d -p $(dirname $SRC_ROOT)`
export PATH=${GOPATH}/bin:/usr/local/go/bin:${PATH}
mkdir -p ${GOPATH}/src/github.com/kubeflow
GO_BUILD_DIR=${GOPATH}/src/github.com/kubeflow/kubebench
ln -s ${SRC_ROOT} ${GO_BUILD_DIR}

cd ${SRC_ROOT}

echo "Build go binaries"
go build github.com/kubeflow/kubebench/controller/cmd/configurator
go build github.com/kubeflow/kubebench/controller/cmd/reporter
echo "Build image ${IMAGE}:${TAG}"
docker build -t ${IMAGE}:${TAG} -f ${DOCKERFILE} .

echo "Push image ${IMAGE}:${TAG}"
gcloud auth activate-service-account --key-file=${GOOGLE_APPLICATION_CREDENTIALS}
gcloud docker -- push "${IMAGE}:${TAG}"

echo "Clean up go build directory"
rm -rf ${GOPATH}

echo "Image ${IMAGE}:${TAG} built successfully"
