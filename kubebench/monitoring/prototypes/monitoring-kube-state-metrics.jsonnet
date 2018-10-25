// @apiVersion 0.1
// @name io.ksonnet.pkg.monitoring-kube-state-metrics
// @description Kube State Metrics installer for monitoring
// @shortDescription Kube State Metrics installer for monitoring
// @param name string Name for the component
// @optionalParam namespace string null Namespace to use for the components

local k = import "k.libsonnet";

local kubeStateMetrics = import "kubebench/monitoring/monitoring-kube-state-metrics.libsonnet";

local namespace = if params.namespace == "null" then env.namespace else params.namespace;

local kubeStateMetricsName = import "param://name";

local kubeStateMetricsInstance = kubeStateMetrics.parts(kubeStateMetricsName, namespace);
kubeStateMetricsInstance.list(kubeStateMetricsInstance.all)
