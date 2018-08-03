import argparse
import logging

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

def deploy_kubeflow(test_case):
  """Deploy Kubeflow."""
  args = parse_args()
  test_dir = test_case.test_suite.test_dir
  src_root_dir = args.src_root_dir
  namespace = args.namespace
  api_client = deploy_utils.create_k8s_client()
  app_dir = deploy_utils.setup_ks_app(
      test_dir, src_root_dir, namespace, args.github_token, api_client)

  # Deploy Kubeflow
  util.run(["ks", "generate", "tf-job-operator", "tf-job-operator",
            "--namespace=" + namespace], cwd=app_dir)
  util.run(["ks", "generate", "argo", "kubeflow-argo", "--name=kubeflow-argo",
            "--namespace=" + namespace], cwd=app_dir)
  apply_command = ["ks", "apply", "default",
                   "-c", "tf-job-operator", "-c", "kubeflow-argo"]
  if args.as_gcloud_user:
    account = deploy_utils.get_gcp_identity()
    logging.info("Impersonate %s", account)
    # If we don't use --as to impersonate the service account then we
    # observe RBAC errors when doing certain operations. The problem appears
    # to be that we end up using the in cluster config (e.g. pod service account)
    # and not the GCP service account which has more privileges.
    apply_command.append("--as=" + account)
  util.run(apply_command, cwd=app_dir)

  # Verify that the TfJob operator is actually deployed.
  tf_job_deployment_name = "tf-job-operator-v1alpha2"
  logging.info("Verifying TfJob controller started.")
  util.wait_for_deployment(api_client, namespace, tf_job_deployment_name)

  # Verify that the Argo operator is deployed.
  argo_deployment_name = "workflow-controller"
  logging.info("Verifying Argo controller started.")
  util.wait_for_deployment(api_client, namespace, argo_deployment_name)

  util.run(["ks", "generate", "nfs-server", "nfs-server", "--name=nfs-server",
            "--namespace=" + namespace], cwd=app_dir)
  apply_command = ["ks", "apply", "default",
                   "-c", "nfs-server"]
  if args.as_gcloud_user:
    account = deploy_utils.get_gcp_identity()
    logging.info("Impersonate %s", account)
    # If we don't use --as to impersonate the service account then we
    # observe RBAC errors when doing certain operations. The problem appears
    # to be that we end up using the in cluster config (e.g. pod service account)
    # and not the GCP service account which has more privileges.
    apply_command.append("--as=" + account)
  util.run(apply_command, cwd=app_dir)
  util.wait_for_deployment(api_client, namespace, "nfs-server")

  nfs_server_ip = deploy_utils.get_nfs_server_ip("nfs-server",namespace)

  util.run(["ks", "generate", "nfs-volume", "nfs-volume", "--name=kubebench-pvc","--nfs_server_ip="+nfs_server_ip,
            "--namespace=" + namespace], cwd=app_dir)
  apply_command = ["ks", "apply", "default",
                   "-c", "nfs-volume"]
  if args.as_gcloud_user:
    account = deploy_utils.get_gcp_identity()
    logging.info("Impersonate %s", account)
    # If we don't use --as to impersonate the service account then we
    # observe RBAC errors when doing certain operations. The problem appears
    # to be that we end up using the in cluster config (e.g. pod service account)
    # and not the GCP service account which has more privileges.
    apply_command.append("--as=" + account)
  util.run(apply_command, cwd=app_dir)
  #util.wait_for_deployment(api_client, namespace, "nfs-volume")
  deploy_utils.copy_job_config(src_root_dir + "/kubeflow/kubebench", namespace)

def main():
  test_case = test_helper.TestCase(
    name='deploy_kubeflow', test_func=deploy_kubeflow)
  test_suite = test_helper.init(
    name='deploy_kubeflow', test_cases=[test_case])
  test_suite.run()

if __name__ == "__main__":
  main()
