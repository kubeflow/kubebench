// @apiVersion 0.1
// @name io.ksonnet.pkg.kubebench-example-tfcnn
// @description kubebench-example-tfcnn
// @shortDescription A simple TFJob to run CNN benchmark
// @param name string Name for the job.
// @optionalParam image string gcr.io/xyhuang-kubeflow/kubebench-example-tfcnn-runner-cpu:v20180909-1 Image
// @optionalParam num_worker number 1 Number of workders
// @optionalParam num_ps number 1 Number of parameter servers
// @optionalParam args string null Other arguments
// @optionalParam batch_size number 2 Batch size
// @optionalParam num_batches number 20 Number of batches to run
// @optionalParam model string resnet50 Model
// @optionalParam variable_update string parameter_server Variable update
// @optionalParam num_gpus number 1 Number of GPUs
// @optionalParam local_parameter_device string cpu Local parameter device
// @optionalParam device string cpu Device
// @optionalParam data_format string NHWC Data format
// @optionalParam forward_only string true Use forward-only mode


local k = import "k.libsonnet";

local name = params.name;
local namespace = env.namespace;
local image = import "param://image";
local numWorker = import "param://num_worker";
local numPs = import "param://num_ps";
local argsStr = import "param://args";
local batchSize = import "param://batch_size";
local numBatches = import "param://num_batches";
local model = import "param://model";
local variableUpdate = import "param://variable_update";
local numGpus = import "param://num_gpus";
local localParameterDevice = import "param://local_parameter_device";
local device = import "param://device";
local dataFormat = import "param://data_format";
local forwardOnly = import "param://forward_only";

local args =
  if argsStr == "null" then
    []
  else
    std.split(argsStr, ",");

local tfjob = {
  apiVersion: "kubeflow.org/v1alpha2",
  kind: "TFJob",
  metadata: {
    name: name,
    namespace: namespace,
  },
  spec: {
    tfReplicaSpecs: {
      Worker: {
        replicas: numWorker,
        template: {
          spec: {
            containers: [
              {
                args: [
                  "--batch_size=" + batchSize,
                  "--num_batches=" + numBatches,
                  "--model=" + model,
                  "--variable_update=" + variableUpdate,
                  "--num_gpus=" + numGpus,
                  "--local_parameter_device=" + localParameterDevice,
                  "--device=" + device,
                  "--data_format=" + dataFormat,
                  "--forward_only=" + forwardOnly,
                ] + args,
                image: image,
                name: "tensorflow",
                workingDir: "/opt/tf-benchmarks/scripts/tf_cnn_benchmarks",
              },
            ],
            restartPolicy: "OnFailure",
          },
        },
      },
    } + if numPs > 0 then {
      Ps: {
        replicas: numPs,
        template: {
          spec: {
            containers: [
              {
                args: [
                  "--batch_size=" + batchSize,
                  "--num_batches=" + numBatches,
                  "--model=" + model,
                  "--variable_update=" + variableUpdate,
                  "--num_gpus=" + numGpus,
                  "--local_parameter_device=" + localParameterDevice,
                  "--device=" + device,
                  "--data_format=" + dataFormat,
                  "--forward_only=" + forwardOnly,
                ] + args,
                image: image,
                name: "tensorflow",
                workingDir: "/opt/tf-benchmarks/scripts/tf_cnn_benchmarks",
              },
            ],
            restartPolicy: "OnFailure",
          },
        },
        tfReplicaType: "PS",
      },
    } else {},
  },
};

k.core.v1.list.new([
  tfjob,
])
