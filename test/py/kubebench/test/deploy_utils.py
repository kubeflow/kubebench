import datetime
import logging
import os
import ssl
import time
import uuid

from kubernetes import client as k8s_client
from kubernetes import config
from kubernetes.client import rest

from kubeflow.testing import util  # pylint: disable=no-name-in-module

def get_gcp_identity():
  google_application_credentials = os.getenv(
      "GOOGLE_APPLICATION_CREDENTIALS", None)
  if google_application_credentials:
    util.run(["gcloud", "auth", "activate-service-account",
              "--key-file=" + google_application_credentials])
  else:
    logging.warning("GOOGLE_APPLICATION_CREDENTIALS not set.")
  identity = util.run(["gcloud", "config", "get-value", "account"])
  logging.info("Current GCP account: %s", identity)
  return identity

def create_k8s_client():
  # We need to load the kube config so that we can have credentials to
  # talk to the APIServer.
  util.load_kube_config(persist_config=False)

  # Create an API client object to talk to the K8s master.
  api_client = k8s_client.ApiClient()

  return api_client

def _setup_test(api_client, run_label):
  """Create the namespace for the test.

  Returns:
    test_dir: The local test directory.
  """

  api = k8s_client.CoreV1Api(api_client)
  namespace = k8s_client.V1Namespace()
  namespace.api_version = "v1"
  namespace.kind = "Namespace"
  namespace.metadata = k8s_client.V1ObjectMeta(
    name=run_label, labels={
      "app": "kubeflow-e2e-test",
    })

  try:
    logging.info("Creating namespace %s", namespace.metadata.name)
    namespace = api.create_namespace(namespace)
    logging.info("Namespace %s created.", namespace.metadata.name)
  except rest.ApiException as e:
    if e.status == 409:
      logging.info("Namespace %s already exists.", namespace.metadata.name)
    else:
      raise

  return namespace

def setup_ks_app(test_dir, src_root_dir, namespace, github_token, api_client):
  """Create a ksonnet app for Kubeflow"""
  util.makedirs(test_dir)

  logging.info("Using test directory: %s", test_dir)

  namespace_name = namespace

  namespace = _setup_test(api_client, namespace_name)
  logging.info("Using namespace: %s", namespace)
  if github_token:
    logging.info("Setting GITHUB_TOKEN to %s.", github_token)
    # Set a GITHUB_TOKEN so that we don't rate limited by GitHub;
    # see: https://github.com/ksonnet/ksonnet/issues/233
    os.environ["GITHUB_TOKEN"] = github_token

  if not os.getenv("GITHUB_TOKEN"):
    logging.warning("GITHUB_TOKEN not set; you will probably hit Github API "
                    "limits.")

  kubeflow_registry = os.path.join(src_root_dir, "kubeflow",
      "kubeflow", "kubeflow")
  kubebench_registry = os.path.join(src_root_dir, "kubeflow",
      "kubebench", "kubebench")

  app_name = "kubeflow-test-" + uuid.uuid4().hex[0:4]
  app_dir = os.path.join(test_dir, app_name)

  # Initialize a ksonnet app.
  util.run(["ks", "init", app_name], cwd=test_dir)

  # Add required registries
  registries = {
    "kubeflow": kubeflow_registry,
    "kubebench": kubebench_registry
  }
  for r in registries:
    util.run(["ks", "registry", "add", r, registries[r]], cwd=app_dir)

  # Install required packages
  packages = ["kubeflow/core", "kubeflow/argo", "kubebench/kubebench-job",
              "kubebench/nfs-server", "kubebench/nfs-volume"]
  for p in packages:
    util.run(["ks", "pkg", "install", p], cwd=app_dir)

  return app_dir

def log_operation_status(operation):
  """A callback to use with wait_for_operation."""
  name = operation.get("name", "")
  status = operation.get("status", "")
  logging.info("Operation %s status %s", name, status)

def wait_for_operation(client,
                       project,
                       op_id,
                       timeout=datetime.timedelta(hours=1),
                       polling_interval=datetime.timedelta(seconds=5),
                       status_callback=log_operation_status):
  """Wait for the specified operation to complete.

  Args:
    client: Client for the API that owns the operation.
    project: project
    op_id: Operation id.
    timeout: A datetime.timedelta expressing the amount of time to wait before
      giving up.
    polling_interval: A datetime.timedelta to represent the amount of time to
      wait between requests polling for the operation status.

  Returns:
    op: The final operation.

  Raises:
    TimeoutError: if we timeout waiting for the operation to complete.
  """
  endtime = datetime.datetime.now() + timeout
  while True:
    try:
      op = client.operations().get(
        project=project, operation=op_id).execute()

      if status_callback:
        status_callback(op)

      status = op.get("status", "")
      # Need to handle other status's
      if status == "DONE":
        return op
    except ssl.SSLError as e:
      logging.error("Ignoring error %s", e)
    if datetime.datetime.now() > endtime:
      raise TimeoutError(
        "Timed out waiting for op: {0} to complete.".format(op_id))
    time.sleep(polling_interval.total_seconds())

  # Linter complains if we don't have a return here even though its unreachable.
  return None

def copy_job_config(src_dir, namespace):

  config.load_kube_config()

  v1 = k8s_client.CoreV1Api()
  nfs_server_pod = None
  ret = v1.list_namespaced_pod(namespace, watch=False)
  for i in ret.items:
    if(i.metadata.labels.get("role") != None) & (i.metadata.labels.get("role") == "nfs-server"):
      nfs_server_pod = i.metadata.name
  if nfs_server_pod is None:
    logging.info("nfs server pod NOT found")
    return 0

  cmd = "kubectl -n " + namespace + " exec " + nfs_server_pod + " -- mkdir -p /exports/config"
  util.run(cmd.split(), cwd=src_dir)

  cmd = "kubectl cp examples/tf_cnn_benchmarks/job_config.yaml " + namespace + \
          "/" + nfs_server_pod + ":/exports/config/job-config.yaml"
  util.run(cmd.split(), cwd=src_dir)

  return 1

def get_nfs_server_ip(name, namespace):

  config.load_kube_config()

  v1 = k8s_client.CoreV1Api()
  server_ip = None
  ret = v1.read_namespaced_service(name, namespace)
  if (ret != None) & (ret.spec.cluster_ip != None):
    server_ip = ret.spec.cluster_ip

  return server_ip

def check_kb_job(job_name, namespace):

  config.load_kube_config()

  crd_api = k8s_client.CustomObjectsApi()
  GROUP = "argoproj.io"
  VERSION = "v1alpha1"
  PLURAL = "workflows"
  res = crd_api.get_namespaced_custom_object(GROUP, VERSION, namespace, PLURAL, job_name)

  if res["status"]["phase"] == "Succeeded":
    logging.info("Job Completed")
    return 1

  cmd = "kubectl get pods -n " + namespace
  util.run(cmd.split(), cwd=app_dir)

  logging.info("Job NOT Completed")
  return 0

def cleanup_kb_job(app_dir, job_name):

  cmd = "ks delete default -c " + job_name
  util.run(cmd.split(), cwd=app_dir)
  cmd = "ks delete default -c nfs-volume"
  util.run(cmd.split(), cwd=app_dir)
  cmd = "ks delete default -c nfs-server"
  util.run(cmd.split(), cwd=app_dir)
  cmd = "ks delete default -c kubeflow-argo"
  util.run(cmd.split(), cwd=app_dir)
  cmd = "ks delete default -c kubeflow-core"
  util.run(cmd.split(), cwd=app_dir)
