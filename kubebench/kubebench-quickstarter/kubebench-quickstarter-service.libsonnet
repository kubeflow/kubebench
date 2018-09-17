local k = import "k.libsonnet";

{
  parts:: {
    nfsDeployment(name, namespace):: {
      apiVersion: "extensions/v1beta1",
      kind: "Deployment",
      metadata: {
        name: name,
        namespace: namespace,
        labels: {
          role: "kubebench-nfs",
        },
      },
      spec: {
        template: {
          metadata: {
            labels: {
              role: "kubebench-nfs",
            },
          },
          spec: {
            volumes: [
              {
                name: "git-repo",
                emptyDir: {},
              },
              {
                name: "kubebench",
                emptyDir: {},
              },
            ],
            initContainers: [
              {
                name: "init-clone-repo",
                image: "alpine/git",
                args: [
                  "clone",
                  "--single-branch",
                  "--",
                  "https://github.com/kubeflow/kubebench.git",
                  "/kubebench/repo",
                ],
                volumeMounts: [
                  {
                    name: "git-repo",
                    mountPath: "/kubebench/repo",
                  },
                ],
              },
              {
                name: "copy-data",
                image: "busybox",
                command: [
                  "sh",
                ],
                args: [
                  "-c",
                  "mkdir -p /mnt/kubebench/config/registry ; " +
                  "mkdir -p /mnt/kubebench/data ; " +
                  "mkdir -p /mnt/kubebench/experiments ; " +
                  "cp -r /kubebench-repo/kubebench /mnt/kubebench/config/registry/ ; " +
                  "cp -r /kubebench-repo/examples/config/* /mnt/kubebench/config",
                ],
                volumeMounts: [
                  {
                    name: "git-repo",
                    mountPath: "/kubebench-repo",
                  },
                  {
                    name: "kubebench",
                    mountPath: "/mnt",
                  },
                ],
              },
            ],
            containers: [
              {
                name: "kubebench-nfs-server",
                image: "k8s.gcr.io/volume-nfs:0.8",
                volumeMounts: [
                  {
                    name: "kubebench",
                    mountPath: "/exports",
                  },
                ],
                ports: [
                  {
                    name: "nfs",
                    containerPort: 2049,
                  },
                  {
                    name: "mountd",
                    containerPort: 20048,
                  },
                  {
                    name: "rpcbind",
                    containerPort: 111,
                  },
                ],
                securityContext: {
                  privileged: true,
                },
              },
              {
                name: "kubebench-nfs-file-server",
                image: "httpd:2.4-alpine",
                command: [
                  "sh",
                ],
                args: [
                  "-c",
                  "rm -f /usr/local/apache2/logs/httpd.pid ; " +
                  'httpd -DFOREGROUND -c "DocumentRoot /usr/local/apache2/htdocs/kubebench"',
                ],
                volumeMounts: [
                  {
                    name: "kubebench",
                    mountPath: "/usr/local/apache2/htdocs",
                  },
                ],
                ports: [
                  {
                    name: "file-server",
                    containerPort: 80,
                  },
                ],
              },
            ],  // containers
          },
        },  // template
      },
    },  // nfsDeployment

    nfsService(name, namespace, serviceType):: {
      kind: "Service",
      apiVersion: "v1",
      metadata: {
        name: name,
        namespace: namespace,
      },
      spec: {
        type: serviceType,
        ports: [
          {
            name: "nfs",
            port: 2049,
          },
          {
            name: "mountd",
            port: 20048,
          },
          {
            name: "rpcbind",
            port: 111,
          },
        ],
        selector: {
          role: "kubebench-nfs",
        },
      },
    },  // nfsService

    nfsFileServerService(name, namespace, serviceType):: {
      kind: "Service",
      apiVersion: "v1",
      metadata: {
        name: name,
        namespace: namespace,
      },
      spec: {
        type: serviceType,
        ports: [
          {
            name: "file-server",
            port: 80,
          },
        ],
        selector: {
          role: "kubebench-nfs",
        },
      },
    },  // nfsFileServerService
  },
}
