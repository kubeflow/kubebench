#!/bin/bash

KB_VERSION=${KB_VERSION:-master}
KF_ENV=${KF_ENV:-default}

ks registry add kubebench github.com/kubeflow/kubebench/tree/${KB_VERSION}/kubebench
ks registry add kubeflow github.com/kubeflow/kubeflow/tree/master/kubeflow

ks pkg install kubeflow/core
ks pkg install kubebench/monitoring

ks generate ambassador ambassador
ks generate monitoring-prometheus-operator monitoring-prometheus-operator
ks generate monitoring-kube-state-metrics monitoring-kube-state-metrics
ks generate monitoring-node-exporter monitoring-node-exporter
ks generate monitoring-prometheus monitoring-prometheus
ks generate monitoring-grafana monitoring-grafana --prometheusName=monitoring-prometheus

ks apply $KF_ENV -c ambassador
ks apply $KF_ENV -c monitoring-prometheus-operator
ks apply $KF_ENV -c monitoring-kube-state-metrics
ks apply $KF_ENV -c monitoring-node-exporter
ks apply $KF_ENV -c monitoring-prometheus
ks apply $KF_ENV -c monitoring-grafana
