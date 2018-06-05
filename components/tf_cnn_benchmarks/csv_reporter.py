import csv
import argparse
import os
import re
import warnings


class LogParser(object):
  """Class LogParser
  Takes path to .txt file as the first argument
  and writes simplified log to .csv file
  Sample execution: python reporter.py path_to_log.txt path_to_table.csv
  Path to .csv file is optional. If not given,
  creates new file in the same folder as the script
  under the name log.csv
  """

  def __init__(self):
    self.columns = [
      'total_images/sec',
      'TensorFlow',
      'Model',
      'Dataset',
      'Mode',
      'SingleSess',
      'Batch_size',
      'Num_batches',
      'Num_epochs',
      'Data_format',
      'Layout_optimizer',
      'Optimizer',
      'Variables',
      'Sync']
    self.output_file = self.parse_command_line().output_file
    self.log_dir = self.parse_command_line().log_dir
    self.log_dict = dict(zip(self.columns, [None] * len(self.columns)))
    self.create_csv()

  def parse_command_line(self):
    self.parser = argparse.ArgumentParser(
      description='Turn logs into csv.')
    self.parser.add_argument('--log-dir', type=str,
                 help='Directory of log files')
    # Optional positional argument
    self.parser.add_argument(
      '--output-file',
      type=str,
      help='Path to csv output file')
    return self.parser.parse_args()

  def create_csv(self):
    output_dir = os.path.dirname(self.output_file)
    if not os.path.exists(output_dir):
      os.makedirs(output_dir)
    if not os.path.isfile(self.output_file):
      warnings.warn("Creating new csv log file")
      with open(self.output_file, 'w+') as csvfile:
        writer = csv.DictWriter(csvfile, fieldnames=self.columns)
        writer.writeheader()

  @classmethod
  def extract_value(cls, line):
    '''
    Args: line - string
    Takes line and checks wheres this line contains
    one of the keys from self.columns
    If contains, stores it and then writes it to .csv file
    '''
    msg = line.strip('\n').split('|')[-1].strip()
    filtered_msg = re.sub('[!@#$\'\""]', '', msg)
    key_value_pair = filtered_msg .strip().split(": ")
    found_key = key_value_pair[0].replace(' ', '_')
    try:
      value = key_value_pair[1]
      value = value.strip(r'\   \\n')
      return found_key, value
    except IndexError:
      return None, None

  def read_logs(self):
    log_file = self.log_dir + "/worker0.log"
    with open(log_file, 'r') as f:
      line = 'init'
      while line:
        line = f.readline()
        key, value = self.extract_value(line)
        if key in self.columns:
          self.log_dict[key] = value

  def write_csv(self):
    with open(self.output_file, 'a') as csvfile:
      writer = csv.DictWriter(csvfile, fieldnames=self.columns)
      writer.writerow(self.log_dict)


if __name__ == "__main__":
  log_parser = LogParser()
  log_parser.read_logs()
  log_parser.write_csv()
