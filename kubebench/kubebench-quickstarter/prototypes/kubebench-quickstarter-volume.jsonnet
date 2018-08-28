// @apiVersion 0.1
// @name io.ksonnet.pkg.kubebench-quickstarter-volume
// @description Kubebench quick-start volume installer
// @shortDescription Kubebench quick-start volume installer
// @param name string Name for the installer.
// @optionalParam namespace string null Namespace to use for the components
// @optionalParam nfsServiceIP string null service type of the nfs file server

local k = import "k.libsonnet";
local nfsVol = import "kubebench/kubebench-quickstarter/kubebench-quickstarter-volume.libsonnet";

local name = import "param://name";
local namespace = if params.namespace == "null" then env.namespace else params.namespace;
local nfsServiceIP = params.nfsServiceIP;

local capacity = "1Gi";
local storageRequest = "1Gi";

std.prune(k.core.v1.list.new([
  nfsVol.parts.nfsPV("kubebench-config-pv", namespace, nfsServiceIP, capacity, "/kubebench/config", "config"),
  nfsVol.parts.nfsPVC("kubebench-config-pvc", namespace, storageRequest, "config"),
  nfsVol.parts.nfsPV("kubebench-data-pv", namespace, nfsServiceIP, capacity, "/kubebench/data", "data"),
  nfsVol.parts.nfsPVC("kubebench-data-pvc", namespace, storageRequest, "data"),
  nfsVol.parts.nfsPV("kubebench-exp-pv", namespace, nfsServiceIP, capacity, "/kubebench/experiments", "experiments"),
  nfsVol.parts.nfsPVC("kubebench-exp-pvc", namespace, storageRequest, "experiments"),
]))
