FROM golang
COPY operator /app/operator
RUN ls /app/operator/
RUN go install /app/operator/
RUN ls /app/operator/
RUN ls /app/
ENTRYPOINT /go/bin/kubebench-operator