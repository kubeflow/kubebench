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

FROM tensorflow/tensorflow:1.8.0-gpu

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    build-essential \
    git

RUN mkdir -p /opt
RUN mkdir -p /kubebench/experiments

RUN git clone -n https://github.com/tensorflow/benchmarks.git /opt/tf-benchmarks
RUN cd /opt/tf-benchmarks; git checkout 3b90c14fb2bf02ca5d27c188aee878663229a0a7

COPY examples/tf-cnn/src/runner.py /opt
RUN chmod u+x /opt/*
WORKDIR /opt/tf-benchmarks/scripts/tf_cnn_benchmarks/
ENTRYPOINT ["python", "/opt/runner.py"]
