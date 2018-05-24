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
                  "name": "step-report",
                  "template": "report",
                }
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
              "args": [
                        "--output-file=" + pvcMount + "/output/" + name + "/manifest.json",
                        "--runner-log-dir=" + pvcMount + "/output/" + name,
                        "--pvc-name=" + pvcName,
                        "--pvc-mount=" + pvcMount,
                      ] + configArgs,
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
                    "path": pvcMount + "/output/" + name + "/manifest.json"
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
            "name": "report",
            "container": {
              "image": reportImage,
              "imagePullPolicy": "IfNotPresent",
              "args": [
                        "--log-dir=" + pvcMount + "/output/" + name,
                      ] + reportArgs,
              "volumeMounts": [
                {
                  "name": "kubebench-volume",
                  "mountPath": pvcMount,
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
