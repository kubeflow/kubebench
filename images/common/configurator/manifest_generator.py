# Copyright 2018 Cisco Systems, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import argparse
import io
import _jsonnet
import yaml


def main():
    parser = argparse.ArgumentParser(description="Convert benchmark configs.")
    parser.add_argument("--config-file", help="config file")
    parser.add_argument("--output-file", help="output file")
    args = parser.parse_args()

    with io.open(args.config_file, "r") as stream:
        params = yaml.load(stream, Loader=yaml.BaseLoader)
    job_type = params["jobType"]
    jsonnet_file = job_type + ".jsonnet"

    json_string = _jsonnet.evaluate_file(
            jsonnet_file, ext_vars=params["jobParams"])
    
    with io.open(args.output_file, "w") as f:
        f.write(json_string)


if __name__ == "__main__":
    main()
