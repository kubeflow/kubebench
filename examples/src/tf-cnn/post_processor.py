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

import json
import os
import sys


def run():
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

  args = sys.argv[1:]
  output_dir = args[0]
  result_dir = args[1]

  result_file = os.path.join(result_dir, "result.json")
  if not os.path.exists(result_dir):
    os.makedirs(result_dir)
  log_file = os.path.join(output_dir, "worker0.log")

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


if __name__ == "__main__":
  run()
