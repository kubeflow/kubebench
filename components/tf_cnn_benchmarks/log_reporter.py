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

import csv
import io
import os
import re
import argparse


def save_csv(arr_results, output_file):
  with io.open(output_file, 'w+') as f:
    wr = csv.writer(f, quoting=csv.QUOTE_ALL)
    wr.writerows(arr_results)


def process_line(line):
  lines = line.split('\t')
  lines[0] = lines[0].split()[-1]
  return lines


def report(input_file, output_file):
  f = io.open(input_file, "r")
  results = f.read().split('\n')
  arr_results = []

  do_process = False
  for res in results:
    if res.find('Done warm up') != -1:
      do_process = True
    if do_process and res.find('---') == -1:
      lines = process_line(res)
      arr_results.append(lines)
  f.close()
  save_csv(arr_results, output_file)


def main():
  parser = argparse.ArgumentParser(description="Report benchmark results.")
  parser.add_argument("--files-path", help="log files dir", action="store", dest="path")
  parser.add_argument("--output-file", help="output file", action="store", dest="outfile")
  args = parser.parse_args()

  for filename in os.listdir(args.path):
    if re.search(r'worker\d+.log', filename):
      report(args.path + '/' + filename, args.outfile)

if __name__ == "__main__":
  main()
