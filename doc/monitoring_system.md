# Guide for Kubebench monitoring system

The monitoring system collects metrics from Kubernetes resources and shows them by using Prometheus and Grafana. The system consists of 5 components: [kube-state-metrics](https://github.com/kubernetes/kube-state-metrics), [node-exporter](https://github.com/prometheus/node_exporter), [Prometheus](https://github.com/prometheus/prometheus), [Grafana](https://github.com/grafana/grafana), [Prometheus-operator](https://github.com/coreos/prometheus-operator).
You need to use [Ambassador](https://github.com/datawire/ambassador) to access Grafana UI.

## Installation

* In your Ksonnet app root, run the followings

```
export KB_VERSION=master
export KF_ENV=default
curl https://raw.githubusercontent.com/kubeflow/kubebench/master/scripts/install_monitoring.sh | bash
```

## View results 

You can check Prometheus service port ```kubectl get svc --namespace ${NAMESPACE}``` to access Prometheus UI.
You can access to Grafana UI using port of the Ambassador and /grafana/ link.
Grafana has pre-installing Kubebench dashboard with information about memory usage, CPU usage, disk IOs. For additional information about using Grafana visit [Grafana](http://docs.grafana.org/) docs.

## Cleanups
* Delete monitoring deployment

```
curl https://raw.githubusercontent.com/kubeflow/kubebench/master/scripts/unistall_monitoring.sh | bash
```
