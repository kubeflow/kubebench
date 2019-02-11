# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
from __future__ import print_function
import json
import os
from kubernetes import client, config
import yaml

def load_yaml_dict(file_name):
  # Step 1. Find MPI Job kubectl get -f /tmp/mpi-job.yaml -o name
  # Kubernetes Python SDK doesn't have api to get from file, parse from yaml instead
  with open(file_name, 'r') as stream:
    try:
      doc = yaml.load(stream)
    except yaml.YAMLError as exc:
      print(exc)

    namespace = doc.get('metadata').get('namespace', 'default')
    job_name = doc.get('metadata').get('name')
    return namespace, job_name

def dump_logs(namespace, job_name, output_dir):
  config.load_incluster_config()
  apiV1 = client.CoreV1Api()

  # Step 2. Find launcher pod
  # kubectl get pods -l mpi_job_name=mpi-job-custom,mpi_role_type=launcher -o name
  label_selector = 'mpi_job_name=' + job_name + "," + "mpi_role_type=launcher"
  pods = apiV1.list_namespaced_pod(namespace, label_selector=label_selector)

  # Step 3. Read pod logs - kubectl logs mpi-job-custom-launcher-rhvt4
  launcher_pod_name = pods.items[0].metadata.name
  api_response = apiV1.read_namespaced_pod_log(launcher_pod_name, namespace, pretty='true')

  log_file = os.path.join(output_dir, launcher_pod_name)
  with open(log_file, 'w') as stream:
    stream.write(api_response)

  return log_file

def parse_logs(log_file, result_file):
  result_keys = [
    'total images/sec',
    'TensorFlow',
    'Model',
    'Dataset',
    'Mode',
    'SingleSess',
    'Batch size',
    'Num batches',
    'Num epochs',
    'Data format',
    'Layout optimizer',
    'Optimizer',
    'Variables',
    'Sync']

  result = {}
  # read logs and populate result dict
  with open(log_file, "r") as f:
    for line in f:
      for rk in result_keys:
        # remove log info and then find key words
        if line.split("|", 4)[-1].lstrip().find(rk + ":") == 0:
          key = rk.lower().replace(" ", "_")
          value = line.split(rk + ":")[-1].strip()
          result[key] = value
          break

  # write result as json file
  with open(result_file, "w") as f:
    json.dump(result, f)


def run():
  config_dir = os.environ.get("KUBEBENCH_EXP_CONFIG_PATH")
  output_dir = os.environ.get("KUBEBENCH_EXP_OUTPUT_PATH")
  result_dir = os.environ.get("KUBEBENCH_EXP_RESULT_PATH")

  config_file = os.path.join(config_dir, "kf-job-manifest.yaml")
  result_file = os.path.join(result_dir, "result.json")

  namespace, job_name = load_yaml_dict(config_file)
  log_file = dump_logs(namespace, job_name, output_dir)
  parse_logs(log_file, result_file)


if __name__ == "__main__":
  run()
