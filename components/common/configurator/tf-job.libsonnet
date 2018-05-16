{
  parts:: {
    tfJobReplica(replicaType, number, args, image, pvcName, pvcMount, numGpus=0)::
      local baseVolume = {
        name: "kubebench-volume",
        persistentVolumeClaim: {
          claimName: pvcName,
        },
      };
      local baseContainer = {
        image: image,
        imagePullPolicy: "IfNotPresent",
        name: "tensorflow",
      };
      local containerArgs = if std.length(args) > 0 then
        {
          args: args,
        }
      else {};
      local resources = if numGpus > 0 then {
        resources: {
          limits: {
            "nvidia.com/gpu": numGpus,
          },
        },
      } else {};
      local volumeMounts = {
        volumeMounts: [
          {
            name: "kubebench-volume",
            mountPath: pvcMount,
          }
        ],
      };
      if number > 0 then
        {
          replicas: number,
          template: {
            spec: {
              volumes: [
                baseVolume,
              ],
              containers: [
                baseContainer + containerArgs + resources + volumeMounts,
              ],
              restartPolicy: "OnFailure",
            },
          },
          tfReplicaType: replicaType,
        }
      else {},

    tfJobTerminationPolicy(replicaName, replicaIndex):: {
      chief: {
        replicaName: replicaName,
        replicaIndex: replicaIndex,
      },
    },

    tfJob(name, namespace, replicas, tp):: {
      apiVersion: "kubeflow.org/v1alpha1",
      kind: "TFJob",
      metadata: {
        name: name,
        namespace: namespace,
      },
      spec: {
        replicaSpecs: replicas,
        terminationPolicy: tp,
      },
    },
  },
}
