// @apiVersion 0.1
// @name io.ksonnet.pkg.kubebench-example-tfcnn
// @description kubebench-example-tfcnn
// @shortDescription A simple TFJob to run CNN benchmark
// @param name string Name for the job.
// @optionalParam image string gcr.io/xyhuang-kubeflow/kubebench-example-tfcnn-runner-cpu:v20180826-1 Image
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
            // TODO (xyhuang): the kubebench volumes are to be automated by the configurator.
            volumes: [
              {
                name: "kubebench-config-volume",
                persistentVolumeClaim: {
                  claimName: "kubebench-config-pvc",
                },
              },
              {
                name: "kubebench-exp-volume",
                persistentVolumeClaim: {
                  claimName: "kubebench-exp-pvc",
                },
              },
            ],
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
                // TODO (xyhuang): the kubebench envs are to be automated by the configurator.
                env: [
                  {
                    name: "KUBEBENCH_CONFIG_ROOT",
                    value: "/kubebench/config",
                  },
                  {
                    name: "KUBEBENCH_EXP_ROOT",
                    value: "/kubebench/experiments",
                  },
                  {
                    name: "KUBEBENCH_DATA_ROOT",
                    value: "/kubebench/data",
                  },
                  {
                    name: "KUBEBENCH_EXP_ID",
                    value: name,
                  },
                  {
                    name: "KUBEBENCH_EXP_DIR",
                    value: "$(KUBEBENCH_EXP_ROOT)/$(KUBEBENCH_EXP_ID)",
                  },
                  {
                    name: "KUBEBENCH_EXP_CONFIG_DIR",
                    value: "$(KUBEBENCH_EXP_DIR)/config",
                  },
                  {
                    name: "KUBEBENCH_EXP_OUTPUT_DIR",
                    value: "$(KUBEBENCH_EXP_DIR)/output",
                  },
                  {
                    name: "KUBEBENCH_EXP_RESULT_DIR",
                    value: "$(KUBEBENCH_EXP_DIR)/result",
                  },
                ],
                // TODO (xyhuang): the kubebench volumes are to be automated by the configurator.
                volumeMounts: [
                  {
                    name: "kubebench-config-volume",
                    mountPath: "/kubebench/config",
                  },
                  {
                    name: "kubebench-exp-volume",
                    mountPath: "/kubebench/experiments",
                  },
                ],
              },
            ],
            restartPolicy: "OnFailure",
          },
        },
      },
      Ps: {
        replicas: numPs,
        template: {
          spec: {
            // TODO (xyhuang): the kubebench volumes are to be automated by the configurator.
            volumes: [
              {
                name: "kubebench-config-volume",
                persistentVolumeClaim: {
                  claimName: "kubebench-config-pvc",
                },
              },
              {
                name: "kubebench-exp-volume",
                persistentVolumeClaim: {
                  claimName: "kubebench-exp-pvc",
                },
              },
            ],
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
                // TODO (xyhuang): the kubebench envs are to be automated by the configurator.
                env: [
                  {
                    name: "KUBEBENCH_CONFIG_ROOT",
                    value: "/kubebench/config",
                  },
                  {
                    name: "KUBEBENCH_EXP_ROOT",
                    value: "/kubebench/experiments",
                  },
                  {
                    name: "KUBEBENCH_DATA_ROOT",
                    value: "/kubebench/data",
                  },
                  {
                    name: "KUBEBENCH_EXP_ID",
                    value: name,
                  },
                  {
                    name: "KUBEBENCH_EXP_DIR",
                    value: "$(KUBEBENCH_EXP_ROOT)/$(KUBEBENCH_EXP_ID)",
                  },
                  {
                    name: "KUBEBENCH_EXP_CONFIG_DIR",
                    value: "$(KUBEBENCH_EXP_DIR)/config",
                  },
                  {
                    name: "KUBEBENCH_EXP_OUTPUT_DIR",
                    value: "$(KUBEBENCH_EXP_DIR)/output",
                  },
                  {
                    name: "KUBEBENCH_EXP_RESULT_DIR",
                    value: "$(KUBEBENCH_EXP_DIR)/result",
                  },
                ],
                // TODO (xyhuang): the kubebench volumes are to be automated by the configurator.
                volumeMounts: [
                  {
                    name: "kubebench-config-volume",
                    mountPath: "/kubebench/config",
                  },
                  {
                    name: "kubebench-exp-volume",
                    mountPath: "/kubebench/experiments",
                  },
                ],
              },
            ],
            restartPolicy: "OnFailure",
          },
        },
        tfReplicaType: "PS",
      },
    },
  },
};

k.core.v1.list.new([
  tfjob,
])
