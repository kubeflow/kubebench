#!/usr/bin/env bash

CURRENT=$(pwd)
ROOT_PACKAGE="github.com/kubeflow/kubebench/contorller/kubebench-operator"
CUSTOM_RESOURCE_NAME=""
CUSTOM_RESOURCE_VERSION

cd $GOPATH/src/k8s.io/code-generator

./generate-groups.sh all "$ROOT_PACKAGE/pkg/client" "$ROOT_PACKAGE/pkg/apis" "$CUSTOM_RESOURCE_NAME:$CUSTOM_RESOURCE_NAME"

cd $CURRENT