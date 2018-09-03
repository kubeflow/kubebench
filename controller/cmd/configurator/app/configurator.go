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
	"path"

	"github.com/ghodss/yaml"
	log "github.com/sirupsen/logrus"

	"github.com/kubeflow/kubebench/controller/pkg/apis/kubebench/v1alpha1"
)

type Configurator struct {
	FileOperator      FileOperatorInterface
	ManifestGenerator ManifestGeneratorInterface
}

func (c *Configurator) Run(config string, manifestOutput string, experimentIDOutput string) error {

	// Read and parse config
	runnerConfigRaw, err := c.FileOperator.ReadConfig(config)
	if err != nil {
		log.Errorf("Failed to read runner config: %s", err)
		return err
	}
	runnerConfig := &v1alpha1.RunnerConfig{}
	if err := yaml.Unmarshal(runnerConfigRaw, runnerConfig); err != nil {
		log.Errorf("Could not parse job runner config; Error: %s\n", err)
		return err
	}

	// Generate experiment ID
	// TODO(xyhuang): add timestamp in experiment ID when ingestion into kf job is implemented
	experimentName := runnerConfig.Metadata.Name
	experimentID := experimentName // + "-" + time.Now().Format("20060102150405")

	// Generate manifest
	manifest, err := c.ManifestGenerator.GenerateManifest(runnerConfig)
	if err != nil {
		log.Errorf("Failed to generate manifest: %s", err)
		return err
	}

	// Write outputs
	outputsMap := map[string][]byte{
		experimentIDOutput: []byte(experimentID),
		manifestOutput:     manifest,
	}
	err = c.FileOperator.WriteOutputs(outputsMap)
	if err != nil {
		log.Errorf("Failed to write outputs: %s", err)
		return err
	}

	// Initialize experiment
	_, configFilename := path.Split(config)
	_, manifestFilename := path.Split(manifestOutput)
	outputsMap = map[string][]byte{
		configFilename:   runnerConfigRaw,
		manifestFilename: manifest,
	}
	err = c.FileOperator.InitExperiment(experimentID, outputsMap)
	if err != nil {
		log.Errorf("Failed to initialize experiment: %s", err)
		return err
	}

	return nil
}
