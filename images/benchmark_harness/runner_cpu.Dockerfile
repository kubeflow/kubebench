# Copyright 2018 Cisco Systems, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM tensorflow/tensorflow:1.7.0

RUN apt-get update
RUN apt-get install -y --no-install-recommends \
    ca-certificates \
    build-essential \
    git

RUN pip install --upgrade google-api-python-client pyyaml paramiko google-cloud

RUN mkdir -p /workspace/git

RUN git clone -n https://github.com/tfboyd/benchmark_harness.git \
    /workspace/git/benchmark_harness

WORKDIR /workspace/git/benchmark_harness

RUN git checkout -b kubebench 01232ff8a1025b9947a319ebb15a031b1283459b

ENTRYPOINT ["python", "-m", "oss_bench.harness.controller"]
