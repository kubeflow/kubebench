import argparse
import logging
from os import path

from kubebench.test import deploy_utils
from kubeflow.testing import test_helper
from kubeflow.testing import util  # pylint: disable=no-name-in-module


def parse_args():
  parser = argparse.ArgumentParser()
  parser.add_argument(
    "--namespace", default=None, type=str, help=("The namespace to use."))
  parser.add_argument(
    "--as_gcloud_user",
    dest="as_gcloud_user",
    action="store_true",
    help=("Impersonate the user corresponding to the gcloud "
          "command with kubectl and ks."))
  parser.add_argument(
    "--no-as_gcloud_user", dest="as_gcloud_user", action="store_false")
  parser.set_defaults(as_gcloud_user=False)
  parser.add_argument(
    "--github_token",
    default=None,
    type=str,
    help=("The GitHub API token to use. This is needed since ksonnet uses the "
          "GitHub API and without it we get rate limited. For more info see: "
          "https://github.com/ksonnet/ksonnet/blob/master/docs"
          "/troubleshooting.md. Can also be set using environment variable "
          "GITHUB_TOKEN."))
  parser.add_argument(
    "--src_root_dir",
    default=None,
    type=str,
    help=("The source directory of all repositories.")
  )

  args, _ = parser.parse_known_args()
  return args

def deploy_kubeflow(test_case): # pylint: disable=unused-argument
  """Deploy Kubeflow."""
  args = parse_args()
  src_root_dir = args.src_root_dir
  namespace = args.namespace
  api_client = deploy_utils.create_k8s_client()

  manifest_repo_dir = path.join(src_root_dir, "kubeflow", "manifests")
  argo_manifest_dir = path.join(manifest_repo_dir, "argo", "base")
  tfoperator_manifest_dir = path.join(manifest_repo_dir, "tf-training",
      "tf-job-operator", "base")

  deploy_utils.setup_test(api_client, namespace)

  apply_args = "-f -"
  if args.as_gcloud_user:
    account = deploy_utils.get_gcp_identity()
    logging.info("Impersonate %s", account)
    # If we don't use --as to impersonate the service account then we
    # observe RBAC errors when doing certain operations. The problem appears
    # to be that we end up using the in cluster config (e.g. pod service account)
    # and not the GCP service account which has more privileges.
    apply_args = " ".join(["--as=" + account, apply_args])

  # Deploy argo
  logging.info("Deploying argo")
  util.run(["kustomize", "edit", "set", "namespace", namespace],
           cwd=argo_manifest_dir)
  util.run(["sh", "-c", "kustomize build | kubectl apply " + apply_args],
           cwd=argo_manifest_dir)

  # Deploy tf-job-operator
  logging.info("Deploying tf-job-operator")
  util.run(["kustomize", "edit", "set", "namespace", namespace],
           cwd=tfoperator_manifest_dir)
  util.run(["sh", "-c", "kustomize build | kubectl apply " + apply_args],
           cwd=tfoperator_manifest_dir)

  # Verify that the TfJob operator is actually deployed.
  tf_job_deployment_name = "tf-job-operator"
  logging.info("Verifying TfJob controller started.")
  util.wait_for_deployment(api_client, namespace, tf_job_deployment_name)

  # Verify that the Argo operator is deployed.
  argo_deployment_name = "workflow-controller"
  logging.info("Verifying Argo controller started.")
  util.wait_for_deployment(api_client, namespace, argo_deployment_name)

  deploy_utils.set_clusterrole(namespace)

def main():
  test_case = test_helper.TestCase(
    name='deploy_kubeflow', test_func=deploy_kubeflow)
  test_suite = test_helper.init(
    name='deploy_kubeflow', test_cases=[test_case])
  test_suite.run()

if __name__ == "__main__":
  main()
