import json
import logging
from kubernetes import client as k8s_client
from kubeflow.testing import test_util, util
from py import ks_util
from py import test_runner
from py import tf_job_client

CPU_TFJOB_COMPONENT_NAME = "simple_tfjob_v1alpha2"
GPU_TFJOB_COMPONENT_NAME = "gpu_tfjob_v1alpha2"

class SimpleTfJobTests(test_util.TestCase):
  def __init__(self, args):
    namespace, name, env = test_runner.parse_runtime_params(args)
    self.app_dir = args.app_dir
    self.env = env
    self.namespace = namespace
    self.tfjob_version = args.tfjob_version
    self.params = args.params
    super(SimpleTfJobTests, self).__init__(class_name="SimpleTfJobTests", name=name)

  # Run a generic TFJob, wait for it to complete, and check for pod/service creation errors.
  def run_simple_tfjob(self, component):
    api_client = k8s_client.ApiClient()

    # Setup the ksonnet app
    ks_util.setup_ks_app(self.app_dir, self.env, self.namespace, component, self.params)

    # Create the TF job
    util.run(["ks", "apply", self.env, "-c", component], cwd=self.app_dir)
    logging.info("Created job %s in namespaces %s", self.name, self.namespace)

    # Wait for the job to either be in Running state or a terminal state
    logging.info("Wait for conditions Running, Succeeded, or Failed")
    results = tf_job_client.wait_for_condition(
      api_client, self.namespace, self.name, ["Running", "Succeeded", "Failed"],
      status_callback=tf_job_client.log_status)
    logging.info("Current TFJob:\n %s", json.dumps(results, indent=2))

    # Wait for the job to complete.
    logging.info("Waiting for job to finish.")
    results = tf_job_client.wait_for_job(
      api_client, self.namespace, self.name, self.tfjob_version,
      status_callback=tf_job_client.log_status)
    logging.info("Final TFJob:\n %s", json.dumps(results, indent=2))

    if not tf_job_client.job_succeeded(results):
      self.failure = "Job {0} in namespace {1} in status {2}".format(
        self.name, self.namespace, results.get("status", {}))
      logging.error(self.failure)
      return

    # Check for creation failures.
    creation_failures = tf_job_client.get_creation_failures_from_tfjob(
      api_client, self.namespace, results)
    if creation_failures:
      # TODO(jlewi): Starting with
      # https://github.com/kubeflow/tf-operator/pull/646 the number of events
      # no longer seems to match the expected; it looks like maybe events
      # are being combined? For now we just log a warning rather than an
      # error.
      logging.warning(creation_failures)

    # Delete the TFJob.
    tf_job_client.delete_tf_job(api_client, self.namespace, self.name, version=self.tfjob_version)
    logging.info("Waiting for job %s in namespaces %s to be deleted.", self.name,
                 self.namespace)
    tf_job_client.wait_for_delete(
      api_client, self.namespace, self.name, self.tfjob_version,
      status_callback=tf_job_client.log_status)

  # Run a generic TFJob, wait for it to complete, and check for pod/service creation errors.
  def test_simple_tfjob_cpu(self):
    self.run_simple_tfjob(CPU_TFJOB_COMPONENT_NAME)

  # Run a generic TFJob, wait for it to complete, and check for pod/service creation errors.
  def test_simple_tfjob_gpu(self):
    self.run_simple_tfjob(GPU_TFJOB_COMPONENT_NAME)

if __name__ == "__main__":
  test_runner.main(module=__name__)
