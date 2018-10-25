local grafanaDashboardKubebenchMonitoring = import "grafana-dashboards/kubebench-monitoring.json";
local k = import "k.libsonnet";

{
  parts(grafanaName, namespace, prometheusName):: {

    //Grafana Datasource Config
    local grafanaDatasource = {
      apiVersion: 1,
      datasources: [
        {
          access: "proxy",
          editable: true,
          name: "prometheus",
          orgId: 1,
          type: "prometheus",
          url: "http://" + prometheusName + "." + namespace + ".svc:9090",
          version: 1,
        },
      ],
    },

    //Grafana Datasource ConfigMap
    local grafanaDatasourceConfigMap = {
      apiVersion: "v1",
      kind: "ConfigMap",
      metadata: {
        name: "grafana-datasources",
        namespace: namespace,
      },
      data: {
        "prometheus.yaml": std.toString(grafanaDatasource),
      },
    },
    grafanaDatasourceConfigMap:: grafanaDatasourceConfigMap,

    //Grafana Dashboard Source Config
    local grafanaDashboardSource = {
      apiVersion: 1,
      providers: [
        {
          folder: "",
          name: "default",
          options: {
            path: "/grafana-dashboard-definitions/default",
          },
          orgId: 1,
          type: "file",
        },
      ],
    },

    //Grafana Dashboard Source ConfigMap
    local grafanaDashboardSourceConfigMap = {
      apiVersion: "v1",
      kind: "ConfigMap",
      metadata: {
        name: "grafana-dashboards",
        namespace: namespace,
      },
      data: {
        "dashboards.yaml": std.toString(grafanaDashboardSource),
      },
    },
    grafanaDashboardSourceConfigMap:: grafanaDashboardSourceConfigMap,

    //Grafana Dasboard Kubebench Monitoring ConfigMap
    local grafanaDashboardKubebenchMonitoringConfigMap = {
      apiVersion: "v1",
      kind: "ConfigMap",
      metadata: {
        name: "grafana-dashboard-kubebench-monitoring",
        namespace: namespace,
      },
      data: {
        "kubebench-monitoring.json": "" + grafanaDashboardKubebenchMonitoring + "",
      },
    },
    grafanaDashboardKubebenchMonitoringConfigMap:: grafanaDashboardKubebenchMonitoringConfigMap,

    //Grafana Deployment
    local deployment = {
      apiVersion: "apps/v1beta2",
      kind: "Deployment",
      metadata: {
        labels: {
          app: "prometheus",
          component: grafanaName,
        },
        name: grafanaName,
        namespace: namespace,
      },
      spec: {
        selector: {
          matchLabels: {
            app: "prometheus",
            component: grafanaName,
          },
        },
        template: {
          metadata: {
            labels: {
              app: "prometheus",
              component: grafanaName,
            },
          },
          spec: {
            containers: [
              {
                image: "grafana/grafana:5.2.1",
                imagePullPolicy: "Always",
                name: grafanaName,
                env: [
                  {
                    name: "GF_SERVER_ROOT_URL",
                    value: "http://localhost:3000/grafana/",
                  },
                ],
                ports: [
                  {
                    containerPort: 3000,
                    name: "http-grafana",
                  },
                ],
                volumeMounts: [
                  {
                    mountPath: "/etc/grafana/provisioning/datasources",
                    name: "grafana-datasources",
                    readOnly: false,
                  },
                  {
                    mountPath: "/etc/grafana/provisioning/dashboards",
                    name: "grafana-dashboards",
                    readOnly: false,
                  },
                  {
                    mountPath: "/grafana-dashboard-definitions/default/kubebench-monitoring",
                    name: "grafana-dashboard-kubebench-monitoring",
                    readOnly: false,
                  },
                ],
              },
            ],
            volumes: [
              {
                name: "grafana-datasources",
                configMap: {
                  name: "grafana-datasources",
                },
              },
              {
                name: "grafana-dashboards",
                configMap: {
                  name: "grafana-dashboards",
                },
              },
              {
                name: "grafana-dashboard-kubebench-monitoring",
                configMap: {
                  name: "grafana-dashboard-kubebench-monitoring",
                },
              },
            ],
          },
        },
      },
    },
    deployment:: deployment,

    //Grafana Service
    local service = {
      apiVersion: "v1",
      kind: "Service",
      metadata: {
        name: grafanaName,
        namespace: namespace,
        annotations: {
          "getambassador.io/config":
            std.join("\n", [
              "---",
              "apiVersion: ambassador/v0",
              "kind:  Mapping",
              "name: grafana-ui-mapping",
              "prefix: /grafana/",
              "service: " + grafanaName + "." + namespace + ":3000",
            ]),
        },  //annotations
      },
      spec: {
        ports: [
          {
            name: "http-grafana",
            port: 3000,
            targetPort: "http-grafana",
          },
        ],
        selector: {
          app: "prometheus",
          component: grafanaName,
        },
      },
    },
    service:: service,

    all:: [
      self.grafanaDatasourceConfigMap,
      self.grafanaDashboardSourceConfigMap,
      self.grafanaDashboardKubebenchMonitoringConfigMap,
      self.deployment,
      self.service,
    ],

    //Create Objects
    list(obj=self.all):: k.core.v1.list.new(obj,),

  },
}
