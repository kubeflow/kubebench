local tfJob = import "tf-job.libsonnet";
// TODO: updatedParams uses the environment namespace if
// the namespace parameter is not explicitly set

local name = std.extVar("name");
local namespace = std.extVar("namespace");
local image = std.extVar("image");
local imageGpu = std.extVar("imageGpu");
local numMasters = std.parseInt(std.extVar("numMasters"));
local numPs = std.parseInt(std.extVar("numPs"));
local numWorkers = std.parseInt(std.extVar("numWorkers"));
local numGpus = std.parseInt(std.extVar("numGpus"));
local pvcName = std.extVar("pvcName");
local pvcMount = std.extVar("pvcMount");
local logDir = std.extVar("logDir");
local argsParam = std.extVar("args");
local args =
  if argsParam == "null" then
    ["--log-dir=" + logDir]
  else
    std.split(argsParam, ",") + ["--log-dir=" + logDir];


local terminationPolicy = if numMasters == 1 then
  tfJob.parts.tfJobTerminationPolicy("MASTER", 0)
else
  tfJob.parts.tfJobTerminationPolicy("WORKER", 0);

local workerSpec = if numGpus > 0 then
  tfJob.parts.tfJobReplica("WORKER", numWorkers, args, imageGpu, numGpus, pvcName, pvcMount)
else
  tfJob.parts.tfJobReplica("WORKER", numWorkers, args, image, pvcName, pvcMount);

std.prune(tfJob.parts.tfJob(
  name,
  namespace,
  [
    tfJob.parts.tfJobReplica("MASTER", numMasters, args, image, pvcName, pvcMount),
    workerSpec,
    tfJob.parts.tfJobReplica("PS", numPs, args, image, pvcName, pvcMount),
  ],
  terminationPolicy
))
