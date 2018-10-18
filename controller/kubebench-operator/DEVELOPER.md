# Kubebench-operator developer guide

Kubebench-operator is a typical controller. A simple one can be found [here](https://github.com/kubernetes/sample-controller).

## Prerequisites

  - GO (define GOPATH, see [guide](https://github.com/golang/go/wiki/SettingGOPATH)
  - dep (dependendcy manager for go) [link](https://github.com/golang/dep)

## How to
In order to run it locally do the following commands:

```bash
go get -u github.com/kubeflow/kubebench/controller/kubebench-operator

cd $GOPATH/github.com/kubeflow/kubebench/controller/kubebench-operator

dep ensure 
```

Now you can build/run the controller locally. For the local use [minikube](https://github.com/kubernetes/minikube) is recommended. 

```bash
minikube start

kubectl apply -f examples/crd.yaml #deploy CRD into k8s cluster

go build . -o kubebench-operator #build operator

go run main.go #or run it locally
```

Client folder is fully generated using codegen. 
Whenever you change ```types.go``` in apis folder or add apis be sure to update your code using ```hack/update-code.sh``` (take a look at env variables needed to run that script correctly)

```bash
./hack/update-code.sh
```

## Controller structure

1. APIs folder is the actual folder containing description of objects of Kind ```KubebenchJob```. ```types.go``` describes Kubebenchjob in terms of Go structures.
2. Client folder is fully generated, it containes clients and interfaces to access kubebench-operator from other programms/controllers.
3. Controller folder contains functions which describe controller logic (addition/deletion/update of the elements).
4. Util folder contains functions used for the conversion of ```KubebenchJob``` to ```Argo Workflow``` using Argo APIs.