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

import argparse
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


def main():
  parser = argparse.ArgumentParser()
  parser.add_argument("--output_dir", type=str)
  # split args for runner process and benchmark subprocess
  runner_args, benchmark_args = parser.parse_known_args(sys.argv[1:])

  tf_config = os.environ.get("TF_CONFIG", '{}')
  tf_config_json = json.loads(tf_config)
  cluster = tf_config_json.get("cluster", {})
  job_name = tf_config_json.get("task", {}).get("type", "")
  task_index = tf_config_json.get("task", {}).get("index", "")

  output_dir = runner_args.output_dir
  if not os.path.exists(output_dir):
    os.makedirs(output_dir)
  jn = job_name if job_name != "" else "worker"
  ti = str(task_index) if task_index != "" else "0"
  output_file = os.path.join(output_dir, jn + ti + ".log")

  command = ["python", "tf_cnn_benchmarks.py"] + benchmark_args
  ps_hosts = ",".join(cluster.get("ps", []))
  worker_hosts = ",".join(cluster.get("worker", []))
  if cluster.get("ps", []) or len(cluster.get("worker", [])) > 1:
    command.append("--job_name=" + job_name)
    command.append("--ps_hosts=" + ps_hosts)
    command.append("--worker_hosts=" + worker_hosts)
    command.append("--task_index=" + str(task_index))

  logging.getLogger().setLevel(logging.INFO)
  logging.basicConfig(level=logging.INFO,
                      filename=output_file,
                      filemode='w',
                      format=('%(levelname)s|%(asctime)s'
                              '|%(pathname)s|%(lineno)d| %(message)s'),
                      datefmt='%Y-%m-%dT%H:%M:%S',
                      )
  logging.info("Runner started.")
  logging.info("Command to run: %s", " ".join(command))

  run_and_stream(command)
  logging.info("Finished: %s", " ".join(command))


if __name__ == "__main__":
  main()
