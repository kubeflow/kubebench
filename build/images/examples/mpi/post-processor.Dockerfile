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

# Build this image from project root folder
# docker build -t ${registry}/${repo}:${tag}  -f build/images/examples/mpi/post-processor.Dockerfile .

FROM python:alpine3.7

LABEL maintainer "seedjeffwan@gmail.com"

RUN apk add --no-cache --update \
    libffi-dev \
    openssl-dev \
    bzip2-dev \
    zlib-dev \
    readline-dev \
    build-base
RUN pip install kubernetes

RUN mkdir /workspace

COPY examples/src/mpi/post_processor.py /workspace

WORKDIR /workspace

ENTRYPOINT ["python", "post_processor.py"]
