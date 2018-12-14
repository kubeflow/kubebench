FROM golang:1.8 AS build

WORKDIR /go/src/app
COPY . .

RUN env CGO_ENABLED=0 go build -a -o /kubebench-dashboard
# ADD . /go/src/github.com/kubeflow/kubebench/dashboard
# RUN go install github.com/kubeflow/kubebench/dashboard

FROM scratch

COPY --from=build /kubebench-dashboard /kubebench-dashboard

EXPOSE 9303
# ENTRYPOINT ./kubebenchjob-controller
ENTRYPOINT ./kubebench-dashboard

# FROM golang:1.8 AS build

# WORKDIR /go/src/app
# COPY . .

# RUN env CGO_ENABLED=0 go build -a -o /kubebenchjob-controller

# FROM scratch

# COPY --from=build /kubebenchjob-controller /kubebenchjob-controller
# ENTRYPOINT ./kubebenchjob-controller