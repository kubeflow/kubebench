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

"""A launcher suitable for invoking tf_cnn_benchmarks using TfJob."""

import logging
import json
import os
import subprocess
import sys


def run_and_stream(cmd):
  logging.info("Running %s", " ".join(cmd))
  process = subprocess.Popen(cmd, stdout=subprocess.PIPE,
                             stderr=subprocess.STDOUT)

  while process.poll() is None:
    process.stdout.flush()
    if process.stderr:
      process.stderr.flush()
    sys.stderr.flush()
    sys.stdout.flush()
    for line in iter(process.stdout.readline, b''):
      process.stdout.flush()
      logging.info(line.strip())

  sys.stderr.flush()
  sys.stdout.flush()
  process.stdout.flush()
  if process.stderr:
    process.stderr.flush()
  for line in iter(process.stdout.readline, b''):
    logging.info(line.strip())

  if process.returncode != 0:
    raise ValueError("cmd: {0} exited with code {1}".format(
      " ".join(cmd), process.returncode))

if __name__ == "__main__":
  tf_config = os.environ.get('TF_CONFIG', '{}')
  tf_config_json = json.loads(tf_config)
  cluster = tf_config_json.get('cluster', {})
  job_name = tf_config_json.get('task', {}).get('type', "")
  task_index = tf_config_json.get('task', {}).get('index', "")

  # pick out the --log-dir arg
  # TODO(xyhuang): this is ugly, needs to be improved
  log_dir = ""
  args = sys.argv[1:]
  for arg in args:
    if arg.startswith("--log-dir="):
      log_dir = arg.lstrip("--log-dir=")
      args.remove(arg)
  command = ["python", "tf_cnn_benchmarks.py"] + args
  ps_hosts = ",".join(cluster.get("ps", []))
  worker_hosts = ",".join(cluster.get("worker", []))
  if len(cluster.get("ps", [])) > 0 or len(cluster.get("worker", [])) > 1:
    command.append("--job_name=" + job_name)
    command.append("--ps_hosts=" + ps_hosts)
    command.append("--worker_hosts=" + worker_hosts)
    command.append("--task_index={0}".format(task_index))

  logging.getLogger().setLevel(logging.INFO)
  logging.basicConfig(level=logging.INFO,
                      filename=log_dir + "/" + job_name + str(task_index) + '.log',
                      filemode='w',
                      format=('%(levelname)s|%(asctime)s'
                              '|%(pathname)s|%(lineno)d| %(message)s'),
                      datefmt='%Y-%m-%dT%H:%M:%S',
                      )
  logging.info("Launcher started.")

  logging.info("Command to run: %s", " ".join(command))
  # with open("/opt/run_benchmarks.sh", "w") as hf:
  #   hf.write("#!/bin/bash\n")
  #   hf.write(" ".join(command))
  #   hf.write("\n")

  run_and_stream(command)
  logging.info("Finished: %s", " ".join(command))
