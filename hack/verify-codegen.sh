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

SCRIPT_DIR=$(dirname $0)
ROOT_PKG="github.com/kubeflow/kubebench"
DIFFROOT="$GOPATH/src/$ROOT_PKG/controller/pkg"
TMP_DIFFROOT=$(mktemp -d tmp.pkg.XXXX)

cleanup() {
  rm -rf $TMP_DIFFROOT
}
trap "cleanup" EXIT SIGINT

echo "Copy ${DIFFROOT}/apis ${DIFFROOT}/client to ${TMP_DIFFROOT}"
cp -a ${DIFFROOT}/apis ${DIFFROOT}/client ${TMP_DIFFROOT}

echo "Run codegen script"
${SCRIPT_DIR}/update-codegen.sh

echo "Check diff"
ret=0
ret=$ret || $(diff -Naupr "${DIFFROOT}/apis" "${TMP_DIFFROOT}/apis")
ret=$ret || $(diff -Naupr "${DIFFROOT}/client" "${TMP_DIFFROOT}/client")

echo "Copy back from ${TMP_DIFFROOT} to ${DIFFROOT}"
cp -a ${TMP_DIFFROOT}/* ${DIFFROOT}

if [[ $ret -eq 0 ]] ; then
  echo "OK"
else
  echo "FAILED"
  echo "Run hack/update-codegen.sh to update the generated code."
  exit 1
fi
