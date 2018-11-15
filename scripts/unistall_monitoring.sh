#!/bin/bash

KF_ENV=${KF_ENV:-default}

ks delete $KF_ENV -c monitoring-grafana
ks delete $KF_ENV -c monitoring-prometheus
ks delete $KF_ENV -c monitoring-kube-state-metrics
ks delete $KF_ENV -c monitoring-node-exporter
ks delete $KF_ENV -c monitoring-prometheus-operator
ks delete $KF_ENV -c ambassador

kubectl delete svc prometheus-operated
