// Tests that each replica has a correctly configured TF RunConfig object.
// Each replica runs a tf-operator-test-server, so a manual exit on the chief
// worker is required for the job to end successfully.
local params = std.extVar("__ksonnet/params").components.estimator_runconfig;

local k = import "k.libsonnet";

local parts(namespace, name, image) = {
  job:: {
    apiVersion: "kubeflow.org/v1alpha2",
    kind: "TFJob",
    metadata: {
      name: name,
      namespace: namespace,
    },
    spec: {
      cleanPodPolicy: "All",
      tfReplicaSpecs: {
        Chief: {
          replicas: 1,
          restartPolicy: "Never",
          template: {
            spec: {
              containers: [
                {
                  name: "tensorflow",
                  image: "gcr.io/kubeflow-images-staging/tf-operator-test-server:v20180904-7d89548b",
                },
              ],
            },
          },
        },
        PS: {
          replicas: 2,
          restartPolicy: "Never",
          template: {
            spec: {
              containers: [
                {
                  name: "tensorflow",
                  image: "gcr.io/kubeflow-images-staging/tf-operator-test-server:v20180904-7d89548b",
                },
              ],
            },
          },
        },
        Worker: {
          replicas: 2,
          restartPolicy: "Never",
          template: {
            spec: {
              containers: [
                {
                  name: "tensorflow",
                  image: "gcr.io/kubeflow-images-staging/tf-operator-test-server:v20180904-7d89548b",
                },
              ],
            },
          },
        },
      },
    },
  },
};

std.prune(k.core.v1.list.new([parts(params.namespace, params.name, params.image).job]))
