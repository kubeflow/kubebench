// @apiVersion 0.1
// @name io.ksonnet.pkg.kubebench-job
// @description A benchmark job on Kubeflow
// @shortDescription A benchmark job on Kubeflow
// @param name string Name to give to each of the components
// @optionalParam namespace string default Namespace
// @optionalParam controller_image string gcr.io/xyhuang-kubeflow/kubebench-controller:v20180909-1 Configurator image
// @optionalParam config_pvc string kubebench-config-pvc Configuration PVC
// @optionalParam data_pvc string null Data PVC
// @optionalParam github_token_secret string null Github token secret
// @optionalParam gcp_credentials_secret string null GCP credentials secret
// @optionalParam experiment_pvc string kubebench-exp-pvc Experiment PVC
// @optionalParam kf_job_config string null Path to the kubeflow job config
// @optionalParam post_processor_image string gcr.io/xyhuang-kubeflow/kubebench-example-tfcnn-postprocessor:v20180909-1 Image of post processor
// @optionalParam post_processor_args string null Arguments of post processor
// @optionalParam reporter_type string csv Type of reporter
// @optionalParam reporter_args string --input-file=result.json,output-file=report.csv Arguments of reporter

local k = import "k.libsonnet";
local kubebench = import "kubebench/kubebench-job/kubebench-job.libsonnet";

local configPvc = import "param://config_pvc";
local controllerImage = import "param://controller_image";
local dataPvc = import "param://data_pvc";
local experimentPvc = import "param://experiment_pvc";
local gcpCredentialsSecret = import "param://gcp_credentials_secret";
local githubTokenSecret = import "param://github_token_secret";
local kfJobConfig = import "param://kf_job_config";
local name = import "param://name";
local namespace = import "param://namespace";
local postProcessorArgsStr = import "param://post_processor_args";
local postProcessorImage = import "param://post_processor_image";
local reporterArgsStr = import "param://reporter_args";
local reporterType = import "param://reporter_type";

local postProcessorArgs =
  if postProcessorArgsStr == "null" then
    []
  else
    std.split(postProcessorArgs, ",");

local reporterArgs =
  if reporterArgsStr == "null" then
    []
  else
    std.split(reporterArgsStr, ",");

std.prune(k.core.v1.list.new([
  kubebench.parts.workflow(name,
                           namespace,
                           controllerImage,
                           configPvc,
                           dataPvc,
                           experimentPvc,
                           githubTokenSecret,
                           gcpCredentialsSecret,
                           kfJobConfig,
                           postProcessorImage,
                           postProcessorArgs,
                           reporterType,
                           reporterArgs),
]))
