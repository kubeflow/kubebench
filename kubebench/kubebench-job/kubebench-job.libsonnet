local k = import "k.libsonnet";

{
  parts:: {

    workflow(name, namespace, configImage, configArgs,
            reportImage, reportArgs, pvcName, pvcMount):: {
      "apiVersion": "argoproj.io/v1alpha1",
      "kind": "Workflow",
      "metadata": {
        "name": name,
        "namespace": namespace,
      },
      "spec": {
        "entrypoint": "kubebench-workflow",
        "volumes": [
          {
            "name": "kubebench-volume",
            "persistentVolumeClaim": {
              "claimName": pvcName,
            },
          },
        ],
        "templates": [
          {
            "name": "kubebench-workflow",
            "steps": [
              [
                {
                  "name": "step-config",
                  "template": "config",
                },
              ],
              [
                {
                  "name": "step-run",
                  "template": "run",
                  "arguments": {
                    "parameters": [
                      {
                        "name": "manifest",
                        "value": "{{steps.step-config.outputs.parameters.manifest}}",
                      },
                    ],
                  },
                },
              ],
              [
                {
                  "name": "step-cleanup",
                  "template": "cleanup",
                  "arguments": {
                    "parameters": [
                      {
                        "name": "manifest",
                        "value": "{{steps.step-config.outputs.parameters.manifest}}",
                      },
                    ],
                  },
                },
              ],
            ],
          },

          {
            "name": "config",
            "container": {
              "image": configImage,
              "imagePullPolicy": "IfNotPresent",
              "args": configArgs,
              "volumeMounts": [
                {
                  "name": "kubebench-volume",
                  "mountPath": pvcMount,
                },
              ],
            },
            "outputs": {
              "parameters": [
                {
                  "name": "manifest",
                  "valueFrom": {
                    "path": pvcMount + "/manifest.json"
                  },
                },
              ],
            },
          },
          {
            "name": "run",
            "resource": {
              "action": "create",
              "successCondition": "status.phase == Done",
              "failureCondition": "status.phase == Failed",
              "manifest": "{{inputs.parameters.manifest}}",
            },
            "inputs": {
              "parameters": [
                {
                  "name": "manifest",
                },
              ],
            },
          },
          {
            "name": "cleanup",
            "resource": {
              "action": "delete",
              "manifest": "{{inputs.parameters.manifest}}",
            },
            "inputs": {
              "parameters": [
                {
                  "name": "manifest",
                },
              ],
            },
          },
        ],
      },
    },
  },
}
