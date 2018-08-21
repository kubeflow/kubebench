FROM golang:alpine AS builder
# The GOPATH in the image is /go.
ADD . /go/src/github.com/kubeflow/kubebench/
# build controller binaries 
WORKDIR /go/src/github.com/kubeflow/kubebench/controller/cmd/configurator
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo
WORKDIR /go/src/github.com/kubeflow/kubebench/controller/cmd/reporter
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo

# pack controller binaries into the image 
FROM alpine:3.7
RUN apk add ca-certificates wget && update-ca-certificates

WORKDIR /app

# install ksonnet
RUN wget https://github.com/ksonnet/ksonnet/releases/download/v0.11.0/ks_0.11.0_linux_amd64.tar.gz 
RUN tar -xzf ks_0.11.0_linux_amd64.tar.gz
ENV PATH="/app/ks_0.11.0_linux_amd64:${PATH}"

#copy built binaries from the builder to the image 
COPY --from=builder /go/src/github.com/kubeflow/kubebench/controller/cmd/configurator/configurator /app/
COPY --from=builder /go/src/github.com/kubeflow/kubebench/controller/cmd/reporter/reporter     /app/

## setup a fake kube config for ks init
COPY --from=builder /go/src/github.com/kubeflow/kubebench/controller/misc/fake_kubeconfig .
ENV KUBECONFIG="/app/fake_kubeconfig"

ENTRYPOINT ["/bin/sh"]
