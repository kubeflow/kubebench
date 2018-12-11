FROM golang

ADD . /go/src/github.com/kubeflow/kubebench/dashboard
RUN go install github.com/kubeflow/kubebench/dashboard

ENTRYPOINT /go/bin/dashboard
EXPOSE 9303
