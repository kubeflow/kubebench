# Guide for Kubebench monitoring system

The monitoring system collects metrics from Kubernetes resources and shows them by using Prometheus and Grafana. The system consists of 4 components: [kube-state-metrics](https://github.com/kubernetes/kube-state-metrics), [node-exporter](https://github.com/prometheus/node_exporter), [Prometheus](https://github.com/prometheus/prometheus), [Grafana](https://github.com/grafana/grafana).

## Installation

* In your Ksonnet app root, run the followings

```
export KB_VERSION=master
export KB_ENV=default
```
* Add registry to your Ksonnet app 

```
ks registry add kubebench github.com/kubeflow/kubebench/tree/${KB_VERSION}/kubebench
```

* Install monitoring package from Kubebench registry

```
ks pkg install kubebench/monitoring
```

* Generate all necessary components

```
ks generate monitoring-kube-state-metrics monitoring-kube-state-metrics
ks generate monitoring-node-exporter monitoring-node-exporter
ks generate monitoring-prometheus monitoring-prometheus
ks generate monitoring-grafana monitoring-grafana --prometheusName=monitoring-prometheus
```

* Apply these components to your environment

```
ks apply ${KB_ENV} -c monitoring-kube-state-metrics
ks apply ${KB_ENV} -c monitoring-node-exporter
ks apply ${KB_ENV} -c monitoring-prometheus
ks apply ${KB_ENV} -c monitoring-grafana
```

## View results 

You can check Prometheus and Grafana services ports ```kubectl get svc --namespace ${NAMESPACE}``` to access UI for Prometheus and Grafana.
Grafana has pre-installing Kubebench dashboard with information about memory usage, CPU usage, disk IOs. For additional information about using Grafana visit [Grafana](http://docs.grafana.org/) docs.

## Cleanups
* Delete monitoring deployment

```
ks delete ${KB_ENV} -c monitoring-kube-state-metrics
ks delete ${KB_ENV} -c monitoring-node-exporter
ks delete ${KB_ENV} -c monitoring-prometheus
ks delete ${KB_ENV} -c monitoring-grafana
```