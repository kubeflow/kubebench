local k = import "k.libsonnet";

{
  parts(nodeExporterName, namespace):: {

    //Node Exporter DaemonSet
    local daemonSet = {
      apiVersion: "apps/v1beta2",
      kind: "DaemonSet",
      metadata: {
        name: nodeExporterName,
        namespace: namespace,
        labels: {
          "k8s-app": "node-exporter",
        },
      },
      spec: {
        selector: {
          matchLabels: {
            "k8s-app": "node-exporter",
          },
        },
        template: {
          metadata: {
            name: nodeExporterName,
            labels: {
              "k8s-app": "node-exporter",
            },
          },
          spec: {
            containers: [
              {
                image: "quay.io/prometheus/node-exporter:v0.16.0",
                imagePullPolicy: "Always",
                name: nodeExporterName,
                ports: [
                  {
                    name: "http-node-exp",
                    containerPort: 9100,
                    hostPort: 9100,
                  },
                ],
              },
            ],
            hostNetwork: true,
            serviceAccountName: nodeExporterName,
          },
        },
      },
    },
    daemonSet:: daemonSet,

    //Node Exporter Service
    local service = {
      apiVersion: "v1",
      kind: "Service",
      metadata: {
        name: nodeExporterName,
        namespace: namespace,
        labels: {
          "k8s-app": "node-exporter",
        },
      },
      spec: {
        ports: [
          {
            name: "http-node-exp",
            port: 9100,
            protocol: "TCP",
          },
        ],
        selector: {
          "k8s-app": "node-exporter",
        },
      },
    },
    service:: service,

    //Node Exporter Service Monitor
    local serviceMonitor = {
      apiVersion: "monitoring.coreos.com/v1",
      kind: "ServiceMonitor",
      metadata: {
        labels: {
          "k8s-app": "node-exporter",
        },
        name: "node-exporter",
        namespace: namespace,
      },
      spec: {
        selector: {
          matchLabels: {
            "k8s-app": "node-exporter",
          },
        },
        endpoints: [
          {
            port: "http-node-exp",
          },
        ],
      },
    },
    serviceMonitor:: serviceMonitor,

    //Node Exporter Service Account
    local serviceAccount = {
      apiVersion: "v1",
      kind: "ServiceAccount",
      metadata: {
        name: nodeExporterName,
        namespace: namespace,
      },
    },
    serviceAccount:: serviceAccount,

    //Node Exporter Role
    local role = {
      apiVersion: "rbac.authorization.k8s.io/v1",
      kind: "Role",
      metadata: {
        name: nodeExporterName,
        namespace: namespace,
      },
      rules: [
        {
          apiGroups: ["authentication.k8s.io"],
          resources: ["tokenreviews"],
          verbs: ["create"],
        },
        {
          apiGroups: ["authorization.k8s.io"],
          resources: ["subjectaccessreviews"],
          verbs: ["create"],
        },
      ],
    },
    role:: role,

    //Node Exporter Role Binding
    local roleBinding = {
      apiVersion: "rbac.authorization.k8s.io/v1",
      kind: "RoleBinding",
      metadata: {
        name: nodeExporterName,
        namespace: namespace,
      },
      roleRef: {
        apiGroup: "rbac.authorization.k8s.io",
        kind: "Role",
        name: nodeExporterName,
      },
      subjects: [
        {
          kind: "ServiceAccount",
          name: nodeExporterName,
          namespace: namespace,
        },
      ],
    },
    roleBinding:: roleBinding,

    all:: [
      self.daemonSet,
      self.service,
      self.serviceMonitor,
      self.serviceAccount,
      self.role,
      self.roleBinding,
    ],

    //Create Objects
    list(obj=self.all):: k.core.v1.list.new(obj,),

  },
}
