#!/bin/bash

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

# This shell script is used to build a cluster and create a namespace from our
# argo workflow


set -o errexit
set -o nounset
set -o pipefail

ZONE="${ZONE}"
PROJECT="${PROJECT}"
CLUSTER_NAME="${CLUSTER_NAME}"
CLUSTER_VERSION="${CLUSTER_VERSION}"

echo "Activating service-account"
gcloud auth activate-service-account --key-file=${GOOGLE_APPLICATION_CREDENTIALS}

echo "Creating GPU cluster"
gcloud --project ${PROJECT} container clusters create ${CLUSTER_NAME} \
    --zone ${ZONE} \
    --cluster-version ${CLUSTER_VERSION}

echo "Configuring kubectl"
gcloud --project ${PROJECT} container clusters get-credentials ${CLUSTER_NAME} \
    --zone ${ZONE}

ACCOUNT=`gcloud config get-value account --quiet`
echo "Grant cluster-admin privileges to account ${ACCOUNT}"
kubectl create clusterrolebinding default-admin --clusterrole=cluster-admin --user=${ACCOUNT}
