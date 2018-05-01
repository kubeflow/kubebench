// @apiVersion 0.1
// @name io.ksonnet.pkg.kubebench-job
// @description A benchmark job on Kubeflow
// @shortDescription A benchmark job on Kubeflow
// @param name string Name to give to each of the components
// @optionalParam namespace string default Namespace
// @optionalParam config_image string null Configurator image
// @optionalParam config_args string null Configurator's arguments
// @optionalParam report_image string null Reporter image
// @optionalParam report_agrs string null Reporter's arguments
// @optionalParam pvc_name string null Persistent volume claim name
// @optionalParam pvc_mount string null Persistent volume claim mount point

local k = import "k.libsonnet";
local kubebench = import "kubebench/kubebench-job/kubebench-job.libsonnet";

local name = import "param://name";
local namespace = import "param://namespace";
local configImage = import "param://config_image";
local configArgsStr = import "param://config_args";
local reportImage = import "param://report_image";
local reportArgsStr = import "param://report_args";
local pvcName = import "param://pvc_name";
local pvcMount = import "param://pvc_mount";

local configArgs =
  if configArgsStr == "null" then
    []
  else
    std.split(configArgsStr, ",");
local reportArgs =
  if reportArgsStr == "null" then
    []
  else
    std.split(reportArgsStr, ",");

std.prune(k.core.v1.list.new([
  kubebench.parts.workflow(name, namespace, configImage, configArgs,
      reportImage, reportArgs, pvcName, pvcMount),
]))
