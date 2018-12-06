// @apiVersion 0.1
// @name io.ksonnet.pkg.monitoring-prometheus
// @description Prometheus installer for monitoring
// @shortDescription Prometheus installer for monitoring
// @param name string Name for the component
// @optionalParam namespace string null Namespace to use for the components

local k = import "k.libsonnet";

local prometheus = import "kubebench/monitoring/monitoring-prometheus.libsonnet";

local namespace = if params.namespace == "null" then env.namespace else params.namespace;

local prometheusName = import "param://name";

local prometheusInstance = prometheus.parts(prometheusName, namespace);
prometheusInstance.list(prometheusInstance.all)
