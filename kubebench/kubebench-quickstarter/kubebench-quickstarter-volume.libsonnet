local k = import "k.libsonnet";

{
  parts:: {
    nfsPV(name, namespace, serverip, capacity, path, label):: {
      apiVersion: "v1",
      kind: "PersistentVolume",
      metadata: {
        name: name,
        namespace: namespace,
        labels: {
          kubebenchVolumeType: label,
        },
      },
      spec: {
        capacity: {
          storage: capacity,
        },
        accessModes: ["ReadWriteMany"],
        nfs: {
          server: serverip,
          path: path,
        },
      },
    },

    nfsPVC(name, namespace, storage_request, label):: {
      apiVersion: "v1",
      kind: "PersistentVolumeClaim",
      metadata: {
        name: name,
        namespace: namespace,
      },
      spec: {
        accessModes: ["ReadWriteMany"],
        storageClassName: "",
        resources: {
          requests: {
            storage: storage_request,
          },
        },
        selector: {
          matchLabels: {
            kubebenchVolumeType: label,
          },
        },
      },
    },
  },
}
