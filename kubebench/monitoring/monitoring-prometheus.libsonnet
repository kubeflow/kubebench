local k = import "k.libsonnet";

{
   parts(prometheusName, namespace, serverIP):: {

     //Port Information
     local nodeExporterPort = 9100,
     local kubeStateMetricsPort = 8080,
     local kubeletMetricsPort = 10250,

     //Config Map with Prometheus Config
     local configMap = {
      apiVersion: "v1",
      kind: "ConfigMap",
      metadata: {
        name: prometheusName,
        namespace: namespace,
      },
      data: {
        "prometheus.yml": (importstr "prometheus.yml") %{
          "node-exporter-ip": serverIP+":"+nodeExporterPort,
          "kube-state-metrics-ip": serverIP+":"+kubeStateMetricsPort,
          "kubelet-metrics-ip": serverIP+":"+kubeletMetricsPort,
        },
      },
     },
     configMap:: configMap,
     
     //Prometheus Deployment
     local deployment = {
      apiVersion: "extensions/v1beta1",
      kind: "Deployment",
      metadata: {
        name: prometheusName,
        namespace: namespace,
      },
      spec: {
        selector: {
          matchLabels: {
            app: "prometheus",
          },
        },
        template: {
          metadata: {
            annotations: {
              "prometheus.io/scrape": "true",
            },
            labels: {
              app: "prometheus",
            },
            name: prometheusName,
            namespace: namespace,
          },
          spec: {
            containers: [
              {
                image: "quay.io/prometheus/prometheus",
                imagePullPolicy: "Always",
                name: prometheusName,
                ports: [
                  {
                    containerPort: 9090,
                    name: "http-prometheus",
                  },
                ],
                volumeMounts: [
                  {
                    mountPath: "/etc/prometheus",
                    name: "config-volume",
                  },
                ],
              },
            ],
            serviceAccountName: prometheusName,
            volumes: [
              {
                name: "config-volume",
                configMap: {
                  name: prometheusName,
                },
              },
            ],
          },
        },
      },
    },
    deployment:: deployment,

    //Prometheus Service
    local service = {
      apiVersion: "v1",
      kind: "Service",
      metadata: {
        labels: {
          name: prometheusName,
        },
        name: prometheusName,
        namespace: namespace,
      },
      spec: {
        ports: [
          {
            name: prometheusName,
            port: 9090,
            protocol: "TCP",
          },
        ],
        selector: {
          app: "prometheus",
        },
        type: "NodePort",
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

    all:: [
      self.configMap,
      self.deployment,
      self.service,
      self.serviceAccount,
      self.clusterRole,
      self.clusterRoleBinding,
    ],

    //Create Objects
    list(obj=self.all)::k.core.v1.list.new(obj,),
    
  },
}
