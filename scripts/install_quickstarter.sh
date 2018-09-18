#!/bin/bash

KB_VERSION=${KB_VERSION:-master}
KB_ENV=${KB_ENV:-default}
K8S_NAMESPACE=${K8S_NAMESPACE:-kubeflow}

ks registry add kubebench github.com/kubeflow/kubebench/tree/${KB_VERSION}/kubebench
ks pkg install kubebench/kubebench-quickstarter@${KB_VERSION}
ks pkg install kubebench/kubebench-examples@${KB_VERSION}

ks generate kubebench-quickstarter-service kubebench-quickstarter-service
ks generate kubebench-quickstarter-volume kubebench-quickstarter-volume

ks apply ${KB_ENV} -c kubebench-quickstarter-service

KB_NFS_IP=`kubectl get svc kubebench-nfs-svc -o=jsonpath={.spec.clusterIP} -n ${K8S_NAMESPACE}`
ks param set kubebench-quickstarter-volume nfsServiceIP ${KB_NFS_IP}
ks apply ${KB_ENV} -c kubebench-quickstarter-volume
