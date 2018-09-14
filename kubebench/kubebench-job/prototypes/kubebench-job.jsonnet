// @apiVersion 0.1
// @name io.ksonnet.pkg.kubebench-job
// @description A benchmark job on Kubeflow
// @shortDescription A benchmark job on Kubeflow
// @param name string Name to give to each of the components
// @optionalParam namespace string null Namespace
// @optionalParam controllerImage string gcr.io/xyhuang-kubeflow/kubebench-controller:v20180913-1 Configurator image
// @optionalParam githubTokenSecret string null Github token secret
// @optionalParam githubTokenSecretKey string null Key of Github token secret
// @optionalParam gcpCredentialsSecret string null GCP credentials secret
// @optionalParam gcpCredentialsSecretKey string null Key of GCP credentials secret
// @optionalParam mainJobKsPrototype string kubebench-example-tfcnn The Ksonnet prototype of the job being benchmarked
// @optionalParam mainJobKsPackage string kubebench-examples The Ksonnet package of the job being benchmarked
// @optionalParam mainJobKsRegistry string github.com/kubeflow/kubebench/tree/master/kubebench The Ksonnet registry of the job being benchmarked
// @optionalParam mainJobConfig string tf-cnn-dummy.yaml Path to the config of the benchmarked job
// @optionalParam experimentConfigPvc string kubebench-config-pvc Configuration PVC
// @optionalParam experimentDataPvc string null Data PVC
// @optionalParam experimentRecordPvc string kubebench-exp-pvc Experiment PVC
// @optionalParam postJobImage string gcr.io/xyhuang-kubeflow/kubebench-example-tfcnn-postprocessor:v20180909-1 Image of post processor
// @optionalParam postJobArgs string null Arguments of post processor
// @optionalParam reporterType string csv Type of reporter
// @optionalParam csvReporterInput string result.json The input of CSV reporter
// @optionalParam csvReporterOutput string report.csv The output of CSV reporter

local k = import "k.libsonnet";
local kubebench = import "kubebench/kubebench-job/kubebench-job.libsonnet";

local name = params.name;
local namespace = if params.namespace == "null" then env.namespace else params.namespace;
local controllerImage = params.controllerImage;
local configPvc = params.experimentConfigPvc;
local dataPvc = params.experimentDataPvc;
local experimentPvc = params.experimentRecordPvc;
local gcpCredentialsSecret = params.gcpCredentialsSecret;
local gcpCredentialsSecretKey = params.gcpCredentialsSecretKey;
local githubTokenSecret = params.githubTokenSecret;
local githubTokenSecretKey = params.githubTokenSecretKey;
local mainJobKsPrototype = params.mainJobKsPrototype;
local mainJobKsPackage = params.mainJobKsPackage;
local mainJobKsRegistry = params.mainJobKsRegistry;
local mainJobConfig = params.mainJobConfig;
local postJobArgsStr = params.postJobArgs;
local postJobImage = params.postJobImage;
local reporterType = params.reporterType;
local csvReporterInput = params.csvReporterInput;
local csvReporterOutput = params.csvReporterOutput;

local postJobArgs =
  if postJobArgsStr == "null" then
    []
  else
    std.split(postJobArgs, ",");

std.prune(k.core.v1.list.new([
  kubebench.parts.workflow(name,
                           namespace,
                           controllerImage,
                           configPvc,
                           dataPvc,
                           experimentPvc,
                           githubTokenSecret,
                           githubTokenSecretKey,
                           gcpCredentialsSecret,
                           gcpCredentialsSecretKey,
                           mainJobKsPrototype,
                           mainJobKsPackage,
                           mainJobKsRegistry,
                           mainJobConfig,
                           postJobImage,
                           postJobArgs,
                           reporterType,
                           csvReporterInput,
                           csvReporterOutput),
]))
