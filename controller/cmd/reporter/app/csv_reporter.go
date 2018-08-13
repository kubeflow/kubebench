// Copyright 2018 Cisco Systems, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package app

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	log "github.com/sirupsen/logrus"
)

type CsvReporterOption struct {
	InputFile  string
	OutputFile string
}

func (opt *CsvReporterOption) AddFlags(args []string) error {
	fs := flag.NewFlagSet("csv", flag.ExitOnError)
	fs.StringVar(&opt.InputFile, "input-file", "", "The input file.")
	fs.StringVar(&opt.OutputFile, "output-file", "", "The output file.")
	err := fs.Parse(args)
	if err != nil {
		return err
	}
	return nil
}

type CsvReporter struct{}

func (cr *CsvReporter) Run(options ReporterOption) error {
	inputFile := options.(*CsvReporterOption).InputFile
	outputFile := options.(*CsvReporterOption).OutputFile

	input, err := cr.readInput(inputFile)
	if err != nil {
		log.Errorf("Failed to read input: %s", err)
		return err
	}

	err = cr.appendResult(outputFile, input)
	if err != nil {
		log.Errorf("Failed to send result: %s", err)
		return err
	}

	return nil
}

func (cr *CsvReporter) readInput(filePath string) (map[string]interface{}, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Errorf("Could not read file: %s. Error: %s", filePath, err)
		return nil, err
	}
	var input map[string]interface{}
	err = json.Unmarshal(data, &input)
	if err != nil {
		log.Errorf("Could not parse input; Error: %s", err)
		return nil, err
	}
	return input, nil
}

func (cr *CsvReporter) appendResult(filePath string, input map[string]interface{}) error {
	dir, _ := path.Split(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Errorf("Failed to create directory: %s", err)
		return err
	}

	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	records := [][]string{}
	record := []string{}
	header := []string{}
	reader := csv.NewReader(f)
	// get header from csv file (if exist)
	header, err = reader.Read()
	if err == io.EOF {
		// if empty csv file, then use input map keys as header
		for key := range input {
			header = append(header, key)
		}
		// add header as an record
		records = append(records, header)
	} else if err != nil {
		return err
	}
	for _, key := range header {
		value, ok := input[key]
		if ok {
			switch value.(type) {
			case string, int, float64, bool:
				record = append(record, fmt.Sprint(value))
			case nil:
				record = append(record, "")
			default:
				log.Warnf("Unknonw value type for key: %s", key)
				record = append(record, "")
			}
		} else {
			log.Warnf("Key not found: %s", key)
			record = append(record, "")
		}
	}
	records = append(records, record)

	writer := csv.NewWriter(f)
	err = writer.WriteAll(records)
	if err != nil {
		return err
	}

	return nil
}
