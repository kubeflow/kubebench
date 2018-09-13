// @apiVersion 0.1
// @name io.ksonnet.pkg.kubebench-quickstarter-service
// @description Kubebench quick-start service installer
// @shortDescription Kubebench quick-start service installer
// @param name string Name for the installer.
// @optionalParam namespace string null Namespace to use for the components
// @optionalParam nfsServiceType string ClusterIP Service type of the nfs service
// @optionalParam nfsFileServerServiceType string ClusterIP Service type of the nfs file server

local k = import "k.libsonnet";
local nfsSvc = import "kubebench/kubebench-quickstarter/kubebench-quickstarter-service.libsonnet";

local name = import "param://name";
local namespace = if params.namespace == "null" then env.namespace else params.namespace;

local nfsServiceType = params.nfsServiceType;
local nfsFileServerServiceType = params.nfsFileServerServiceType;

std.prune(k.core.v1.list.new([
  nfsSvc.parts.nfsDeployment("kubebench-nfs-deploy", namespace),
  nfsSvc.parts.nfsService("kubebench-nfs-svc", namespace, nfsServiceType),
  nfsSvc.parts.nfsFileServerService("kubebench-nfs-file-server-svc", namespace, nfsFileServerServiceType),
]))
