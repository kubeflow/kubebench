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
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/kubeflow/kubebench/controller/pkg/constants"
	"github.com/kubeflow/kubebench/controller/pkg/workflowagent/configurator/common"
	"github.com/kubeflow/kubebench/controller/pkg/workflowagent/configurator/manifestgen"
	"github.com/kubeflow/kubebench/controller/pkg/workflowagent/configurator/manifestmod"
)

type Configurator struct {
	fileOperator FileOperatorInterface
}

func NewConfigurator(fo FileOperatorInterface) *Configurator {
	configurator := &Configurator{fileOperator: fo}
	return configurator
}

func (c *Configurator) Run(opt *AppOption) error {

	var configSpec common.ConfiguratorInput

	if err := json.Unmarshal([]byte(opt.InputParams), &configSpec); err != nil {
		log.Errorf("Cannot unmarshal config spec: %s", opt.InputParams)
		return err
	}

	// Generate manifest from the manifest source
	genSpec := configSpec.ManifestGenSpec
	manifest, err := manifestgen.NewManifestGenerator(genSpec).GenerateManifest()
	if err != nil {
		log.Errorf("Failed to generate manifest: %s", err)
		return err
	}

	// Modify the manifest as specified by the kubebench job
	modSpec := configSpec.ManifestModSpec
	// If the namespace in the ManifestModSpec is empty valued, replace it with
	// configurator's own namespace which is same as the workflow. This is for
	// cases when the ManifestModSpec is created before namespace is set.
	if err := replaceEmptyNamespace(modSpec); err != nil {
		log.Errorf("Failed to replace empty namespace: %s", err)
		return err
	}
	modifiedManifest, err := manifestmod.NewManifestModifier(modSpec).ModifyManifest(manifest)
	if err != nil {
		log.Errorf("Failed to modify manifest: %s", err)
		return err
	}

	// Initialize experiment
	experimentID := os.Getenv(constants.ExpIDEnvName)
	err = c.fileOperator.InitExperiment(experimentID, map[string][]byte{})
	if err != nil {
		log.Errorf("Failed to initialize experiment: %s", err)
		return err
	}
	// Write outputs
	outputsMap := map[string][]byte{
		opt.OutputFile: modifiedManifest,
	}
	err = c.fileOperator.WriteOutputs(outputsMap)
	if err != nil {
		log.Errorf("Failed to write outputs: %s", err)
		return err
	}

	return nil
}

// replaceEmptyNamespace replaces an empty-valued namespace with local namespace
func replaceEmptyNamespace(modSpec *common.ManifestModSpec) error {
	if modSpec.Namespace == "" {
		data, err := ioutil.ReadFile(constants.NamespaceFile)
		if err != nil {
			return err
		}
		modSpec.Namespace = strings.Trim(string(data), " \n")
	}
	return nil
}
