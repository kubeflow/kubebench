{
  components: {
    kubebench: {
      name: "kubebench-job-example",
      namespace: "default",
      config_image: "gcr.io/xyhuang-kubeflow/kubebench-configurator:0.0.1",
      config_args: "--config-file=/configurator/job_config.yaml,--output-file=/configurator/manifest.json",
      report_image: "null",
      report_args: "null",
      pvc_mount: "/configurator",
      pvc_name: "kubebench-configurator-pvc",
    },
  },
}
