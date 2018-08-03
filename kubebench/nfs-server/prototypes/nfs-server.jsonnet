// @apiVersion 0.1
// @name io.ksonnet.pkg.nfs-server
// @description NFS server
// @shortDescription Create a nfs server
// @param name string Name for the nfs server.
// @optionalParam namespace string null Namespace to use for the components


local k = import "k.libsonnet";
local nfs = import "ciscoai/nfs-server/nfs-server.libsonnet";

// updatedParams uses the environment namespace if
// the namespace parameter is not explicitly set
local updatedParams = params {
  namespace: if params.namespace == "null" then env.namespace else params.namespace,
};


local name = import "param://name";
local namespace = updatedParams.namespace;

std.prune(k.core.v1.list.new([
  nfs.parts.nfsdeployment(name,namespace),
  nfs.parts.nfsservice(name, namespace)
]))
