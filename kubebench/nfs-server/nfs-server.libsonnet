local k = import "k.libsonnet";

{
 parts:: {
     nfsdeployment(name, namespace):: {
      apiVersion: "extensions/v1beta1",
      kind: "Deployment",
      metadata: {
        name: name,
        namespace: namespace,
        labels: { 
             role: "nfs-server"    
        },
      },
      spec: {
        template: {
          metadata: {
            labels: {
                 role: "nfs-server"
            }
          },
          spec: {
            containers: [ {
                 name: "nfs-server",
                 image: "k8s.gcr.io/volume-nfs:0.8",
                 ports: [ {
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
                 }
                ],
                securityContext: {
                     privileged: true 
                }
              }
            ],
          },
        },
      },
    },  // nfsdeployment

    nfsservice(name, namespace):: {
        kind: "Service",
        apiVersion: "v1",
        metadata: {
            name: name,
            namespace: namespace,
        },
        spec: {
            ports: [{
                name: "nfs", 
                port: 2049
                },
                {
                name: "mountd", 
                port: 20048
                },
                {
                name: "rpcbind", 
                port: 111
                }
            ],
            selector: {
                role: "nfs-server"
            }
        }
    }
  }
}
