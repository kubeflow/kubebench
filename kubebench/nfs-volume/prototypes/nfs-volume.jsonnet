// @apiVersion 0.1
// @name io.ksonnet.pkg.nfs-volume
// @description NFS Persistent volume
// @shortDescription Create a NFS persistent volume claim
// @param name string Name for the nfs volume.
// @param nfs_server_ip string Cluster Ip address of nfs-server service
// @optionalParam namespace string null Namespace to use for the components
// @optionalParam capacity string 1Gi Total capacity of the NFS persistent volume
// @optionalParam mountpath string / NFS mount point
// @optionalParam storage_request string 1Gi Total storage requested from the persistent volume


local k = import "k.libsonnet";
local nfs = import "ciscoai/nfs-volume/nfs-volume.libsonnet";

// updatedParams uses the environment namespace if
// the namespace parameter is not explicitly set
local updatedParams = params {
  namespace: if params.namespace == "null" then env.namespace else params.namespace,
};


local name = import "param://name";
local namespace = updatedParams.namespace;


local nfs_server_ip = import "param://nfs_server_ip";
local capacity = import "param://capacity";
local path = import "param://mountpath";
local storage_request = import "param://storage_request";

std.prune(k.core.v1.list.new([
  nfs.parts.nfsPV(name, namespace, nfs_server_ip, capacity, path),
  nfs.parts.nfsPVC(name, namespace, storage_request)
]))

