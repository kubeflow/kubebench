// Copyright 2019 The Kubeflow Authors
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
	"bytes"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/apimachinery/pkg/util/yaml"

	kbjobv1alpha2 "github.com/kubeflow/kubebench/controller/pkg/apis/kubebenchjob/v1alpha2"
	kbworkflow "github.com/kubeflow/kubebench/controller/pkg/workflow"
)

// Run executes the workflow builder
func Run(opt *AppOption) error {
	// read config
	configDecoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(defaultConfig)), 128)
	if opt.Config != "" {
		configFileReader, err := os.Open(opt.Config)
		if err != nil {
			log.Errorf("Failed to open config file: %s", opt.Config)
			return err
		}
		defer configFileReader.Close()
		configDecoder = yaml.NewYAMLOrJSONDecoder(configFileReader, 128)
	}
	config := &kbjobv1alpha2.KubebenchConfig{}
	if !opt.NoDefaultConfig {
		if err := configDecoder.Decode(config); err != nil {
			return err
		}
	}

	// read manifest
	manifestFile := opt.Manifest
	manifestFileReader, err := os.Open(manifestFile)
	if err != nil {
		log.Errorf("Failed to open manifest file: %s", manifestFile)
		return err
	}
	defer manifestFileReader.Close()
	manifestDecoder := yaml.NewYAMLOrJSONDecoder(manifestFileReader, 128)
	kbjob := &kbjobv1alpha2.KubebenchJob{}
	if err := manifestDecoder.Decode(kbjob); err != nil {
		return err
	}

	// create output writer
	outputFile := opt.Output
	var outputWriter io.Writer
	if outputFile == "" {
		outputWriter = os.Stdout
	} else {
		fileWriter, err := os.Create(outputFile)
		if err != nil {
			log.Errorf("Failed to create output file: %s", outputFile)
			return err
		}
		defer fileWriter.Close()
		outputWriter = fileWriter
	}

	workflow, err := kbworkflow.BuildWorkflow(kbjob, config, false)
	if err != nil {
		log.Errorf("Failed to build Workflow: %s", err)
		return err
	}

	err = json.NewYAMLSerializer(json.DefaultMetaFactory, nil, nil).Encode(workflow, outputWriter)
	if err != nil {
		log.Errorf("Failed to write output: %s", err)
		return err
	}

	return nil
}
