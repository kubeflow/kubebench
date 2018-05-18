import csv
import argparse
import os
import re
import warnings


"""
Class LogParser

Takes path to .txt file as the first argument
and writes simplified log to .csv file


Sample execution: python reporter.py path_to_log.txt path_to_table.csv

Path to .csv file is optional. If not given,
creates new file in the same folder as the script
under the name log.csv

"""


class LogParser:
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
            'Devices',
            'Data_format',
            'Layout_optimizer',
            'Optimizer',
            'Variables',
            'Sync']
        self.path_to_dataframe = self.parse_command_line().path_to_dataframe
        self.path_to_file = self.parse_command_line().filename
        self.log_dict = dict(zip(self.columns, [None] * len(self.columns)))
        self.load_csv_log()

    def parse_command_line(self):
        self.parser = argparse.ArgumentParser(
            description='Parser for filenames')
        self.parser.add_argument('filename', type=str,
                                 help='A required path to txt log file')
        # Optional positional argument
        self.parser.add_argument(
            'path_to_dataframe',
            type=str,
            nargs='?',
            default='log.csv',
            help='An optional path to csv file')
        return self.parser.parse_args()

    def load_csv_log(self):
        if not os.path.isfile(self.path_to_dataframe):
            warnings.warn("Creating new csv log file")
            with open(self.path_to_dataframe, 'w+') as csvfile:
                writer = csv.DictWriter(csvfile, fieldnames=self.columns)
                writer.writeheader()

    def extract_value(self, line):
        '''
        Args: line - string

        Takes line and checks wheres this line contains
        one of the keys from self.columns
        If contains, stores it and then writes it to .csv file
        '''
        msg = line.strip('\n').split('|')[-1].strip()
        filtered_msg = re.sub('[!@#$\'\""]', '', msg[1:])
        key_value_pair = filtered_msg .strip().split(": ")
        found_key = key_value_pair[0].replace(' ', '_')
        try:
            value = key_value_pair[1]
            value = value.strip('\   \\n')
            return found_key, value
        except IndexError:
            return None, None

    def read_logs(self):
        with open(self.path_to_file, 'r') as f:
            line = 'init'
            while line:
                line = f.readline()
                key, value = self.extract_value(line)
                if key in self.columns:
                    self.log_dict[key] = value

    def write_log(self):
        with open(self.path_to_dataframe, 'a') as csvfile:
            writer = csv.DictWriter(csvfile, fieldnames=self.columns)
            writer.writerow(self.log_dict)


if __name__ == "__main__":
    log_parser = LogParser()
    log_parser.read_logs()
    log_parser.write_log()
