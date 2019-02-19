{
  // TODO(https://github.com/ksonnet/ksonnet/issues/222): Taking namespace as an argument is a work around for the fact that ksonnet
  // doesn't support automatically piping in the namespace from the environment to prototypes.

  // convert a list of two items into a map representing an environment variable
  // TODO(jlewi): Should we move this into kubeflow/core/util.libsonnet
  listToMap:: function(v)
    {
      name: v[0],
      value: v[1],
    },

  // Function to turn comma separated list of prow environment variables into a dictionary.
  parseEnv:: function(v)
    local pieces = std.split(v, ",");
    if v != "" && std.length(pieces) > 0 then
      std.map(
        function(i) $.listToMap(std.split(i, "=")),
        std.split(v, ",")
      )
    else [],

  // default parameters.
  defaultParams:: {
    project:: "kubeflow-ci",
    zone:: "us-east1-d",
    registry:: "gcr.io/kubeflow-ci",
    versionTag:: null,
    gcpCredentialsSecretName:: "kubeflow-testing-credentials",
  },

  parts(namespace, name, overrides):: {
    // Workflow to run the e2e test.
    e2e(prow_env, bucket):
      local params = $.defaultParams + overrides;

      // mountPath is the directory where the volume to store the test data
      // should be mounted.
      local mountPath = "/mnt/" + "test-data-volume";
      // testDir is the root directory for all data for a particular test run.
      local testDir = mountPath + "/" + name;
      // outputDir is the directory to sync to GCS to contain the output for this job.
      local outputDir = testDir + "/output";
      local artifactsDir = outputDir + "/artifacts";
      // Source directory where all repos should be checked out
      local srcRootDir = testDir + "/src";
      local srcDir = srcRootDir + "/kubeflow/kubebench";
      // The directory containing the py scripts for testing
      local srcTestPyDir = srcDir + "/test/py";
      // The directory within the kubeflow_testing submodule containing
      // py scripts to use.
      local srcKubeTestPyDir = srcRootDir + "/kubeflow/testing/py";
      local image = "gcr.io/kubeflow-ci/test-worker";
      // The name of the NFS volume claim to use for test files.
      // local nfsVolumeClaim = "kubeflow-testing";
      local nfsVolumeClaim = "nfs-external";
      // The name to use for the volume to use to contain test data.
      local dataVolume = "kubeflow-test-volume";
      // GCP information
      local zone = params.zone;
      local project = params.project;
      // The name of test cluster
      local clusterName = "kubebench-e2e-" + std.substr(name, std.length(name) - 4, 4);
      // The Kubernetes version of test cluster
      local clusterVersion = "1.11";
      // Container build information
      local registry = params.registry;
      local versionTag = if params.versionTag != null then
        params.versionTag else "";

      {
        // Build an Argo template to execute a particular command.
        // step_name: Name for the template
        // command: List to pass as the container command.
        buildTemplate(step_name, command, envVars=[], sidecars=[], workingDir=null, kubeConfig="config"):: {
          name: step_name,
          container: {
            command: command,
            image: image,
            [if workingDir != null then "workingDir"]: workingDir,
            env: [
              {
                // Add the source directories to the python path.
                name: "PYTHONPATH",
                value: srcTestPyDir + ":" + srcKubeTestPyDir,
              },
              {
                name: "GOOGLE_APPLICATION_CREDENTIALS",
                value: "/secret/gcp-credentials/key.json",
              },
              {
                name: "GITHUB_TOKEN",
                valueFrom: {
                  secretKeyRef: {
                    name: "github-token",
                    key: "github_token",
                  },
                },
              },
              {
                name: "ZONE",
                value: zone,
              },
              {
                name: "PROJECT",
                value: project,
              },
              {
                name: "CLUSTER_NAME",
                value: clusterName,
              },
              {
                name: "CLUSTER_VERSION",
                value: clusterVersion,
              },
              {
                // We use a directory in our NFS share to store our kube config.
                // This way we can configure it on a single step and reuse it on subsequent steps.
                name: "KUBECONFIG",
                value: testDir + "/.kube/" + kubeConfig,
              },
              {
                name: "REGISTRY",
                value: registry,
              },
            ] + prow_env + envVars,
            volumeMounts: [
              {
                name: dataVolume,
                mountPath: mountPath,
              },
              {
                name: "github-token",
                mountPath: "/secret/github-token",
              },
              {
                name: "gcp-credentials",
                mountPath: "/secret/gcp-credentials",
              },
            ],
          },
          sidecars: sidecars,
        },  // buildTemplate

        apiVersion: "argoproj.io/v1alpha1",
        kind: "Workflow",
        metadata: {
          name: name,
          namespace: namespace,
        },
        // TODO(jlewi): Use OnExit to run cleanup steps.
        spec: {
          entrypoint: "e2e",
          volumes: [
            {
              name: "github-token",
              secret: {
                secretName: "github-token",
              },
            },
            {
              name: "gcp-credentials",
              secret: {
                secretName: "kubeflow-testing-credentials",
              },
            },
            {
              name: dataVolume,
              persistentVolumeClaim: {
                claimName: nfsVolumeClaim,
              },
            },
          ],  // volumes
          // onExit specifies the template that should always run when the workflow completes.
          onExit: "exit-handler",
          templates: [
            {
              name: "e2e",
              steps: [
                [
                  {
                    name: "checkout",
                    template: "checkout",
                  },
                ],
                [
                  {
                    name: "create-pr-symlink",
                    template: "create-pr-symlink",
                  },
                  {
                    name: "py-test",
                    template: "py-test",
                  },
                  {
                    name: "py-lint",
                    template: "py-lint",
                  },
                  {
                    name: "test-jsonnet-formatting",
                    template: "test-jsonnet-formatting",
                  },
                ],
                [
                  {
                    name: "setup-cluster",
                    template: "setup-cluster",
                  },
                ],
                [
                  {
                    name: "build-kubebench-controller",
                    template: "build-kubebench-controller",
                  },
                  {
                    name: "build-kubebench-examples",
                    template: "build-kubebench-examples",
                  },
                  {
                    name: "build-kubebench-operator",
                    template: "build-kubebench-operator",
                  },
                  {
                    name: "build-kubebench-dashboard",
                    template: "build-kubebench-dashboard",
                  },
                ],
                [
                  {
                    name: "deploy-kubeflow",
                    template: "deploy-kubeflow",
                  },
                ],
                [
                  {
                    name: "wait-for-kubeflow-deployment",
                    template: "wait-for-kubeflow-deployment",
                  },
                ],
                // [
                //   {
                //     name: "test-kubebench-job",
                //     template: "test-kubebench-job",
                //   },
                // ],
              ],
            },
            {
              name: "exit-handler",
              steps: [
                [
                  {
                    name: "teardown-cluster",
                    template: "teardown-cluster",
                  },
                ],
                [
                  {
                    name: "copy-artifacts",
                    template: "copy-artifacts",
                  },
                ],
                [
                  {
                    name: "delete-test-dir",
                    template: "delete-test-dir",
                  },
                ],
              ],
            },
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate(
              "checkout",
              ["/usr/local/bin/checkout.sh", srcRootDir],
              envVars=[{
                name: "EXTRA_REPOS",
                value: "kubeflow/kubeflow@HEAD",
              }],
            ),  // checkout
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("create-pr-symlink", [
              "python",
              "-m",
              "kubeflow.testing.prow_artifacts",
              "--artifacts_dir=" + outputDir,
              "create_pr_symlink",
              "--bucket=" + bucket,
            ]),  // create-pr-symlink
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("copy-artifacts", [
              "python",
              "-m",
              "kubeflow.testing.prow_artifacts",
              "--artifacts_dir=" + outputDir,
              "copy_artifacts",
              "--bucket=" + bucket,
            ]),  // copy-artifacts
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("py-test", [
              "python",
              "-m",
              "kubeflow.testing.test_py_checks",
              "--artifacts_dir=" + artifactsDir,
              "--src_dir=" + srcDir,
            ]),  // py test
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("py-lint", [
              "python",
              "-m",
              "kubeflow.testing.test_py_lint",
              "--artifacts_dir=" + artifactsDir,
              "--src_dir=" + srcDir,
            ]),  // py lint
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("test-jsonnet-formatting", [
              "python",
              "-m",
              "kubeflow.testing.test_jsonnet_formatting",
              "--artifacts_dir=" + artifactsDir,
              "--src_dir=" + srcDir,
            ]),  // test-jsonnet-formatting
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("setup-cluster", [
              srcDir + "/test/scripts/create_cluster.sh",
            ]),  // setup cluster
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("teardown-cluster", [
              srcDir + "/test/scripts/delete_cluster.sh",
            ]),  // teardown cluster
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate(
              "build-kubebench-controller",
              [
                srcDir + "/build/images/controller/build_image.sh",
                srcDir,
                srcDir + "/build/images/controller/Dockerfile",
                "kubebench-controller",
                versionTag,
              ],
              workingDir=srcDir,
            ),  // build-kubebench-controller
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate(
              "build-kubebench-operator",
              [
                srcDir + "/build/images/kubebench-operator/build_image.sh",
                srcDir,
                srcDir + "/build/images/kubebench-operator/Dockerfile",
                "kubebench-operator",
                versionTag,
              ],
              workingDir=srcDir,
            ),  // build-kubebench-operator
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate(
              "build-kubebench-dashboard",
              [
                srcDir + "/build/images/dashboard/build_image.sh",
                srcDir,
                srcDir + "/build/images/dashboard/Dockerfile",
                "kubebench-dashboard",
                versionTag,
              ],
              workingDir=srcDir,
            ),  // build-kubebench-dashboard
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate(
              "build-kubebench-examples",
              [
                srcDir + "/build/images/examples/build_image.sh",
                srcDir,
                srcDir + "/build/images/examples",
                "kubebench-example",
                versionTag,
              ],
              workingDir=srcDir,
            ),  // build-kubebench-examples
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("deploy-kubeflow", [
              "python",
              "-m",
              "kubebench.test.deploy_kubeflow",
              "--test_dir=" + testDir,
              "--src_root_dir=" + srcRootDir,
              "--namespace=" + namespace,
              "--as_gcloud_user",
            ]),  // deploy kubeflow
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("wait-for-kubeflow-deployment", [
              "python",
              "-m",
              "kubebench.test.wait_for_deployment",
              "--cluster=" + clusterName,
              "--project=" + project,
              "--zone=" + zone,
              "--timeout=5",
            ]),  // deploy kubeflow
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("test-kubebench-job", [
              "python",
              "-m",
              "kubebench.test.test_kubebench_job",
              "--test_dir=" + testDir,
              "--src_root_dir=" + srcRootDir,
              "--namespace=" + namespace,
              "--as_gcloud_user",
            ]),  // test kubebench job
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("delete-test-dir", [
              "bash",
              "-c",
              "rm -rf " + testDir,
            ]),  // delete test dir
          ],  // templates
        },
      },  // e2e
  },  // parts
}
