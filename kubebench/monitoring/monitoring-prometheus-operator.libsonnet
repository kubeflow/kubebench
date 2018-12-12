local k = import "k.libsonnet";
local prometheusCRD = import "prometheus-crd.libsonnet";
local serviceMonitorCRD = import "service-monitor-crd.libsonnet";

{
  parts(prometheusOperatorName, namespace):: {

    //Prometheus CRD
    prometheusCRD:: prometheusCRD,

    //Service Monitor CRD
    serviceMonitorCRD:: serviceMonitorCRD,

    //Prometheus Operator Deployment
    local deployment = {
      apiVersion: "apps/v1beta2",
      kind: "Deployment",
      metadata: {
        labels: {
          "k8s-app": "prometheus-operator",
        },
        name: prometheusOperatorName,
        namespace: namespace,
      },
      spec: {
        selector: {
          matchLabels: {
            "k8s-app": "prometheus-operator",
          },
        },
        template: {
          metadata: {
            labels: {
              "k8s-app": "prometheus-operator",
            },
          },
          spec: {
            containers: [
              {
                image: "quay.io/coreos/prometheus-operator:v0.25.0",
                args: [
                  "--kubelet-service=kube-system/kubelet",
                  "--logtostderr=true",
                  "--config-reloader-image=quay.io/coreos/configmap-reload:v0.0.1",
                  "--prometheus-config-reloader=quay.io/coreos/prometheus-config-reloader:v0.25.0",
                ],
                imagePullPolicy: "Always",
                name: prometheusOperatorName,
                ports: [
                  {
                    containerPort: 8080,
                    name: "http-prom-oper",
                  },
                ],
              },
            ],
            serviceAccountName: prometheusOperatorName,
          },
        },
      },
    },
    deployment:: deployment,


    //Prometheus Operator Service
    local service = {
      apiVersion: "v1",
      kind: "Service",
      metadata: {
        labels: {
          "k8s-app": "prometheus-operator",
        },
        name: prometheusOperatorName,
        namespace: namespace,
      },
      spec: {
        ports: [
          {
            name: "http-prom-oper",
            port: 8080,
          },
        ],
        selector: {
          "k8s-app": "prometheus-operator",
        },
      },
    },
    service:: service,

    //Prometheus Service Account
    local serviceAccount = {
      apiVersion: "v1",
      kind: "ServiceAccount",
      metadata: {
        name: prometheusOperatorName,
        namespace: namespace,
      },
    },
    serviceAccount:: serviceAccount,

    //Prometheus Cluster Role
    local clusterRole = {
      apiVersion: "rbac.authorization.k8s.io/v1beta1",
      kind: "ClusterRole",
      metadata: {
        name: prometheusOperatorName,
      },
      rules: [
        {
          apiGroups: [
            "apiextensions.k8s.io",
          ],
          resources: [
            "customresourcedefinitions",
          ],
          verbs: [
            "*",
          ],
        },
        {
          apiGroups: [
            "monitoring.coreos.com",
          ],
          resources: [
            "alertmanagers",
            "prometheuses",
            "prometheuses/finalizers",
            "alertmanagers/finalizers",
            "servicemonitors",
            "prometheusrules",
          ],
          verbs: [
            "*",
          ],
        },
        {
          apiGroups: [
            "apps",
          ],
          resources: [
            "statefulsets",
          ],
          verbs: [
            "*",
          ],
        },
        {
          apiGroups: [
            "",
          ],
          resources: [
            "configmaps",
            "secrets",
          ],
          verbs: [
            "*",
          ],
        },
        {
          apiGroups: [
            "",
          ],
          resources: [
            "pods",
          ],
          verbs: [
            "list",
            "delete",
          ],
        },
        {
          apiGroups: [
            "",
          ],
          resources: [
            "services",
            "endpoints",
          ],
          verbs: [
            "get",
            "create",
            "update",
          ],
        },
        {
          apiGroups: [
            "",
          ],
          resources: [
            "nodes",
          ],
          verbs: [
            "list",
            "watch",
          ],
        },
        {
          apiGroups: [
            "",
          ],
          resources: [
            "namespaces",
          ],
          verbs: [
            "get",
            "list",
            "watch",
          ],
        },
      ],
    },
    clusterRole:: clusterRole,

    //Prometheus Cluster Role
    local clusterRoleBinding = {
      apiVersion: "rbac.authorization.k8s.io/v1beta1",
      kind: "ClusterRoleBinding",
      metadata: {
        name: prometheusOperatorName,
      },
      roleRef: {
        apiGroup: "rbac.authorization.k8s.io",
        kind: "ClusterRole",
        name: prometheusOperatorName,
      },
      subjects: [
        {
          kind: "ServiceAccount",
          name: prometheusOperatorName,
          namespace: namespace,
        },
      ],
    },
    clusterRoleBinding:: clusterRoleBinding,

    all:: [
      self.prometheusCRD,
      self.serviceMonitorCRD,
      self.deployment,
      self.service,
      self.serviceAccount,
      self.clusterRole,
      self.clusterRoleBinding,
    ],

    //Create Objects
    list(obj=self.all):: k.core.v1.list.new(obj,),

  },
}
