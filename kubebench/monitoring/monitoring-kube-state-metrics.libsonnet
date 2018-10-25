local k = import "k.libsonnet";

{
  parts(kubeStateMetricsName, namespace):: {

    //Kube State Metrics Deployment
    local deployment = {
      apiVersion: "apps/v1beta2",
      kind: "Deployment",
      metadata: {
        name: kubeStateMetricsName,
        namespace: namespace,
        labels: {
          app: "prometheus",
          component: kubeStateMetricsName,
        },
      },
      spec: {
        selector: {
          matchLabels: {
            app: "prometheus",
            component: kubeStateMetricsName,
          },
        },
        template: {
          metadata: {
            name: kubeStateMetricsName,
            namespace: namespace,
            labels: {
              app: "prometheus",
              component: kubeStateMetricsName,
            },
          },
          spec: {
            containers: [
              {
                name: kubeStateMetricsName,
                image: "quay.io/coreos/kube-state-metrics:v1.4.0",
                imagePullPolicy: "Always",
                args: [
                  "--namespace=" + namespace,
                ],
                ports: [
                  {
                    name: "http-kube-st-m",
                    containerPort: 8080,
                    hostPort: 8080,
                  },
                ],
              },
            ],
            hostNetwork: true,
            serviceAccountName: kubeStateMetricsName,
          },
        },
      },

    },
    deployment:: deployment,

    //Kube State Metrics Service
    local service = {
      apiVersion: "v1",
      kind: "Service",
      metadata: {
        annotations: {
          "prometheus.io/scrape": "true",
        },
        name: kubeStateMetricsName,
        namespace: namespace,
        labels: {
          app: "prometheus",
          component: kubeStateMetricsName,
        },
      },
      spec: {
        ports: [
          {
            name: "http-kube-st-m",
            port: 8080,
            protocol: "TCP",
          },
        ],
        selector: {
          app: "prometheus",
          component: kubeStateMetricsName,
        },
      },
    },
    service:: service,

    //Kube State Metrics Serivce Account
    local serviceAccount = {
      apiVersion: "v1",
      kind: "ServiceAccount",
      metadata: {
        name: kubeStateMetricsName,
        namespace: namespace,
      },
    },
    serviceAccount:: serviceAccount,


    //Kube State Metrics Cluster Role
    local clusterRole = {
      apiVersion: "rbac.authorization.k8s.io/v1",
      kind: "ClusterRole",
      metadata: {
        name: kubeStateMetricsName,
      },
      rules: [
        {
          apiGroups: [
            "",
          ],
          resources: [
            "configmaps",
            "secrets",
            "nodes",
            "pods",
            "services",
            "resourcequotas",
            "replicationcontrollers",
            "limitranges",
            "persistentvolumeclaims",
            "persistentvolumes",
            "namespaces",
            "endpoints",
          ],
          verbs: [
            "list",
            "watch",
          ],
        },
        {
          apiGroups: [
            "extensions",
          ],
          resources: [
            "daemonsets",
            "deployments",
            "replicasets",
          ],
          verbs: [
            "list",
            "watch",
          ],
        },
        {
          apiGroups: [
            "apps",
          ],
          resources: [
            "statefulsets",
            "daemonsets",
            "deployments",
            "replicasets",
          ],
          verbs: [
            "list",
            "watch",
          ],
        },
        {
          apiGroups: [
            "batch",
          ],
          resources: [
            "cronjobs",
            "jobs",
          ],
          verbs: [
            "list",
            "watch",
          ],
        },
        {
          apiGroups: [
            "autoscaling",
          ],
          resources: [
            "horizontalpodautoscalers",
          ],
          verbs: [
            "list",
            "watch",
          ],
        },
        {
          apiGroups: [
            "authentication.k8s.io",
          ],
          resources: [
            "tokenreviews",
          ],
          verbs: [
            "create",
          ],
        },
        {
          apiGroups: [
            "authorization.k8s.io",
          ],
          resources: [
            "subjectaccessreviews",
          ],
          verbs: [
            "create",
          ],
        },
      ],
    },
    clusterRole:: clusterRole,

    //Kube State Metrics Cluster Role Binding
    local clusterRoleBinding = {
      apiVersion: "rbac.authorization.k8s.io/v1",
      kind: "ClusterRoleBinding",
      metadata: {
        name: kubeStateMetricsName,
      },
      roleRef: {
        apiGroup: "rbac.authorization.k8s.io",
        kind: "ClusterRole",
        name: kubeStateMetricsName,
      },
      subjects: [
        {
          kind: "ServiceAccount",
          name: kubeStateMetricsName,
          namespace: namespace,
        },
      ],
    },
    clusterRoleBinding:: clusterRoleBinding,

    //Kube State Metrics Role
    local role = {
      apiVersion: "rbac.authorization.k8s.io/v1",
      kind: "Role",
      metadata: {
        name: kubeStateMetricsName,
        namespace: namespace,
      },
      rules: [
        {
          apiGroups: [""],
          resources: ["pods"],
          verbs: ["get"],
        },
        {
          apiGroups: ["extensions"],
          resourceNames: [kubeStateMetricsName],
          resources: ["deployments"],
          verbs: [
            "get",
            "update",
          ],
        },
        {
          apiGroups: ["apps"],
          resourceNames: [kubeStateMetricsName],
          resources: ["deployments"],
          verbs: [
            "get",
            "update",
          ],
        },
      ],
    },
    role:: role,

    //Kube State Metrics Role Binding
    local roleBinding = {
      apiVersion: "rbac.authorization.k8s.io/v1",
      kind: "RoleBinding",
      metadata: {
        name: kubeStateMetricsName,
        namespace: namespace,
      },
      roleRef: {
        apiGroup: "rbac.authorization.k8s.io",
        kind: "Role",
        name: kubeStateMetricsName,
      },
      subjects: [
        {
          kind: "ServiceAccount",
          name: kubeStateMetricsName,
          namespace: namespace,
        },
      ],
    },
    roleBinding:: roleBinding,

    all:: [
      self.deployment,
      self.service,
      self.serviceAccount,
      self.clusterRole,
      self.clusterRoleBinding,
      self.role,
      self.roleBinding,
    ],

    //Create Objects
    list(obj=self.all):: k.core.v1.list.new(obj,),

  },
}
