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

FROM alpine:3.7
RUN apk --update add ca-certificates wget && update-ca-certificates

WORKDIR /app

# install ksonnet
RUN wget https://github.com/ksonnet/ksonnet/releases/download/v0.11.0/ks_0.11.0_linux_amd64.tar.gz 
RUN tar -xzf ks_0.11.0_linux_amd64.tar.gz
ENV PATH="/app/ks_0.11.0_linux_amd64:${PATH}"
ENV PATH="/app:${PATH}"

#copy built binaries from the builder to the image 
COPY configurator /app/configurator
COPY resource-manager /app/resource-manager

## setup a fake kube config for ks init
COPY controller/misc/fake_kubeconfig /app/kubeconfig
ENV KUBECONFIG="/app/kubeconfig"
