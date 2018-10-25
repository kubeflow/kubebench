// @apiVersion 0.1
// @name io.ksonnet.pkg.monitoring-node-exporter
// @description Node Exporter installer for monitoring
// @shortDescription Node Exporter installer for monitoring
// @param name string Name for the component
// @optionalParam namespace string null Namespace to use for the components

local k = import "k.libsonnet";

local nodeExporter = import "kubebench/monitoring/monitoring-node-exporter.libsonnet";

local namespace = if params.namespace == "null" then env.namespace else params.namespace;

local nodeExporterName = import "param://name";

local nodeExporterInstance = nodeExporter.parts(nodeExporterName, namespace);
nodeExporterInstance.list(nodeExporterInstance.all)
