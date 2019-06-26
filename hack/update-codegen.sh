#!/bin/bash

# Copyright 2019 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

export GO111MODULE=on

SCRIPT_DIR=$(dirname $0)
ROOT_PKG="github.com/kubeflow/kubebench"
CODEGEN_PKG="k8s.io/code-generator"
CODEGEN_VERSION=$(awk '/k8s.io\/code-generator/ {print $2}' $GOPATH/src/$ROOT_PKG/go.mod)
CODEGEN_PATH="$GOPATH/pkg/mod/$CODEGEN_PKG@$CODEGEN_VERSION"
CRD_NAME="kubebenchjob"
CRD_VERSIONS="v1alpha1,v1alpha2"

TMP_DIR=$(mktemp -d)
echo "Creating temporary directory $TMP_DIR"

mkdir -p $TMP_DIR/$CODEGEN_PKG
cp -r $CODEGEN_PATH/* $TMP_DIR/$CODEGEN_PKG/
chmod -R +w $TMP_DIR/$CODEGEN_PKG

bash $TMP_DIR/$CODEGEN_PKG/generate-groups.sh all \
  "$ROOT_PKG/controller/pkg/client" \
  "$ROOT_PKG/controller/pkg/apis" \
  "$CRD_NAME:$CRD_VERSIONS" \
  --go-header-file "$SCRIPT_DIR/boilerplate/boilerplate.go.txt"

echo "Deleting temporary directory $TMP_DIR"
rm -rf $TMP_DIR
