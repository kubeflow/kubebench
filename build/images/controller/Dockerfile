FROM golang:alpine AS builder
# The GOPATH in the image is /go.
RUN apk add --no-cache git
RUN go get github.com/golang/dep/cmd/dep
ADD . /go/src/github.com/kubeflow/kubebench/
WORKDIR /go/src/github.com/kubeflow/kubebench/
# Install library dependencies
RUN dep ensure -vendor-only

# build controller binaries 
WORKDIR /go/src/github.com/kubeflow/kubebench/controller/cmd/configurator
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo
WORKDIR /go/src/github.com/kubeflow/kubebench/controller/cmd/reporter
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo

# pack controller binaries 
FROM alpine:3.7
RUN apk add ca-certificates wget && update-ca-certificates

WORKDIR /app

RUN wget https://github.com/ksonnet/ksonnet/releases/download/v0.11.0/ks_0.11.0_linux_amd64.tar.gz 
RUN tar -xzf ks_0.11.0_linux_amd64.tar.gz
ENV PATH="/app/ks_0.11.0_linux_amd64:${PATH}"

COPY --from=builder /go/src/github.com/kubeflow/kubebench/controller/cmd/configurator/configurator /app/
COPY --from=builder /go/src/github.com/kubeflow/kubebench/controller/cmd/reporter/reporter     /app/
COPY --from=builder /go/src/github.com/kubeflow/kubebench/controller/misc/fake_kubeconfig .

ENV KUBECONFIG="/app/fake_kubeconfig"

ENTRYPOINT ["/bin/sh"]
