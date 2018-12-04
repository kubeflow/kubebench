local k = import "k.libsonnet";

{
  parts(prometheusName, namespace):: {

    //Prometheus Custom CRD
    local prometheusCustomCRD = {
      apiVersion: "monitoring.coreos.com/v1",
      kind: "Prometheus",
      metadata: {
        name: prometheusName,
        namespace: namespace,
        labels: {
          prometheus: prometheusName,
        },
      },
      spec: {
        baseImage: "quay.io/prometheus/prometheus",
        replicas: 1,
        serviceAccountName: prometheusName,
        serviceMonitorNamespaceSelector: {},
        serviceMonitorSelector: {},
        version: "v2.5.0",
      },
    },
    prometheusCustomCRD:: prometheusCustomCRD,

    //Prometheus Service
    local service = {
      apiVersion: "v1",
      kind: "Service",
      metadata: {
        labels: {
          prometheus: prometheusName,
        },
        name: prometheusName,
        namespace: namespace,
      },
      spec: {
        ports: [
          {
            name: "http-prometheus",
            port: 9090,
          },
        ],
        selector: {
          app: "prometheus",
          prometheus: prometheusName,
        },
      },
    },
    service:: service,

    //Prometheus Service Account
    local serviceAccount = {
      apiVersion: "v1",
      kind: "ServiceAccount",
      metadata: {
        name: prometheusName,
        namespace: namespace,
      },
    },
    serviceAccount:: serviceAccount,

    //Prometheus Cluster Role
    local clusterRole = {
      apiVersion: "rbac.authorization.k8s.io/v1beta1",
      kind: "ClusterRole",
      metadata: {
        name: prometheusName,
      },
      rules: [
        {
          apiGroups: [
            "",
          ],
          resources: [
            "nodes",
            "nodes/proxy",
            "services",
            "endpoints",
            "pods",
            "nodes/metrics",
          ],
          verbs: [
            "get",
            "list",
            "watch",
          ],
        },
        {
          apiGroups: [
            "extensions",
          ],
          resources: [
            "ingresses",
          ],
          verbs: [
            "get",
            "list",
            "watch",
          ],
        },
        {
          nonResourceURLs: [
            "/metrics",
          ],
          verbs: [
            "get",
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
        name: prometheusName,
      },
      roleRef: {
        apiGroup: "rbac.authorization.k8s.io",
        kind: "ClusterRole",
        name: prometheusName,
      },
      subjects: [
        {
          kind: "ServiceAccount",
          name: prometheusName,
          namespace: namespace,
        },
      ],
    },
    clusterRoleBinding:: clusterRoleBinding,

    //Service Monitor Kubelet Metrics
    local serviceMonitorKubelet = {
      apiVersion: "monitoring.coreos.com/v1",
      kind: "ServiceMonitor",
      metadata: {
        labels: {
          "k8s-app": "kubelet",
        },
        name: "kubelet",
        namespace: namespace,
      },
      spec: {
        endpoints: [
          {
            bearerTokenFile: "/var/run/secrets/kubernetes.io/serviceaccount/token",
            interval: "30s",
            port: "https-metrics",
            scheme: "https",
            tlsConfig: {
              insecureSkipVerify: true,
            },
          },
          {
            bearerTokenFile: "/var/run/secrets/kubernetes.io/serviceaccount/token",
            interval: "30s",
            path: "/metrics/cadvisor",
            port: "https-metrics",
            scheme: "https",
            tlsConfig: {
              insecureSkipVerify: true,
            },
          },
        ],
        namespaceSelector: {
          matchNames: [
            "kube-system",
          ],
        },
        selector: {
          matchLabels: {
            "k8s-app": "kubelet",
          },
        },
      },
    },
    serviceMonitorKubelet:: serviceMonitorKubelet,

    all:: [
      self.prometheusCustomCRD,
      self.service,
      self.serviceAccount,
      self.clusterRole,
      self.clusterRoleBinding,
      self.serviceMonitorKubelet,
    ],

    //Create Objects
    list(obj=self.all):: k.core.v1.list.new(obj,),

  },
}
