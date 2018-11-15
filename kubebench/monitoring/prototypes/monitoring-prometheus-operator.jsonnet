// @apiVersion 0.1
// @name io.ksonnet.pkg.monitoring-prometheus-operator
// @description Prometheus operator installer for monitoring
// @shortDescription Prometheus operator installer for monitoring
// @param name string Name for the component
// @optionalParam namespace string null Namespace to use for the components

local k = import "k.libsonnet";

local prometheusOperator = import "kubebench/monitoring/monitoring-prometheus-operator.libsonnet";

local namespace = if params.namespace == "null" then env.namespace else params.namespace;

local prometheusOperatorName = import "param://name";

local prometheusOperatorInstance = prometheusOperator.parts(prometheusOperatorName, namespace);
prometheusOperatorInstance.list(prometheusOperatorInstance.all)
