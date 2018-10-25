// @apiVersion 0.1
// @name io.ksonnet.pkg.monitoring-grafana
// @description Grafana installer for monitoring
// @shortDescription Grafana installer for monitoring
// @param name string Name for the component
// @param prometheusName string Name for Prometheus component
// @optionalParam namespace string null Namespace to use for the components

local k = import "k.libsonnet";

local grafana = import "kubebench/monitoring/monitoring-grafana.libsonnet";

local namespace = if params.namespace == "null" then env.namespace else params.namespace;

local grafanaName = import "param://name";

local prometheusName = import "param://prometheusName";

local grafanaInstance = grafana.parts(grafanaName, namespace, prometheusName);
grafanaInstance.list(grafanaInstance.all)
