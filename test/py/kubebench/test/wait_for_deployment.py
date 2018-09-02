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

"""Wait for kubeflow deployment.

Right now, it only checks for the presence of tfjobs crd.
More things can be added incrementally.

python -m testing.wait_for_deployment --cluster=kubeflow-testing \\
    --project=kubeflow-ci --zone=us-east1-d --timeout=3

"""

from __future__ import print_function

import argparse
import datetime
import logging
import subprocess
import time

from kubeflow.testing import test_helper, util


def parse_args():
  parser = argparse.ArgumentParser()
  parser.add_argument(
    "--cluster",
    default="",
    type=str,
    help="Cluster Name")
  parser.add_argument(
    "--zone",
    default="",
    type=str,
    help="Zone Name")
  parser.add_argument(
    "--project",
    default="",
    type=str,
    help="Project ID")
  parser.add_argument(
    "--timeout",
    default=5,
    type=int,
    help="Timeout in minutes")
  args, _ = parser.parse_known_args()
  return args

def wait_for_resource(resource, end_time):
  while True:
    if datetime.datetime.now() > end_time:
      raise RuntimeError("Timed out waiting for " + resource)
    try:
      if 'error' not in util.run(["kubectl", "get", resource]).lower():
        logging.info("Found %s.", resource)
        break
    except subprocess.CalledProcessError:
      logging.info("Could not find %s. Sleeping for 10 seconds..", resource)
      time.sleep(10)

def test_wait_for_deployment(test_case): # pylint: disable=redefined-outer-name,unused-argument
  args = parse_args()
  util.maybe_activate_service_account()
  util.run([
      "gcloud", "container", "clusters", "get-credentials", args.cluster,
      "--zone=" + args.zone, "--project=" + args.project])
  end_time = datetime.datetime.now() + datetime.timedelta(0, args.timeout*60)
  wait_for_resource("crd/tfjobs.kubeflow.org", end_time)
  wait_for_resource("crd/workflows.argoproj.io", end_time)
  logging.info("Found all resources successfully")

if __name__ == "__main__":
  test_case = test_helper.TestCase(
    name="wait_for_deployment", test_func=test_wait_for_deployment)
  test_suite = test_helper.init(
    name="test_wait_for_deployment", test_cases=[test_case])
  test_suite.run()
