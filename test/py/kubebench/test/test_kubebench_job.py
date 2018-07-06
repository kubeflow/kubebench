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

def run_smoke_test(test_case):
  """Run a smoke test."""
  args = parse_args()
  test_dir = test_case.test_suite.test_dir
  src_root_dir = args.src_root_dir
  namespace = args.namespace
  api_client = deploy_utils.create_k8s_client()
  app_dir = deploy_utils.setup_ks_app(
      test_dir, src_root_dir, namespace, args.github_token, api_client)


  # Deploy Kubebench
  util.run(["ks", "generate", "kubebench-job", "smoke-test-job",
            "--name=smoke-test-job", "--namespace=" + namespace], cwd=app_dir)

  apply_command = ["ks", "apply", "default", "-c", "smoke-test-job"]
  if args.as_gcloud_user:
    account = deploy_utils.get_gcp_identity()
    logging.info("Impersonate %s", account)
    # If we don't use --as to impersonate the service account then we
    # observe RBAC errors when doing certain operations. The problem appears
    # to be that we end up using the in cluster config (e.g. pod service account)
    # and not the GCP service account which has more privileges.
    apply_command.append("--as=" + account)
  # TODO(xyhuang): Currently a place holder so the job is not actually run.
  # A real smoke test job needs to be added and run here.
  # util.run(apply_command, cwd=app_dir)


def main():
  test_case = test_helper.TestCase(
    name='run_smoke_test', test_func=run_smoke_test)
  test_suite = test_helper.init(
    name='test_kubebench_job', test_cases=[test_case])
  test_suite.run()

if __name__ == "__main__":
  main()
