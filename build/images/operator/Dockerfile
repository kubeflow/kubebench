FROM golang
COPY operator /app/operator
RUN go install 	/app/operator

RUN /app/operator
ENTRYPOINT /go/bin/kubebench-operator