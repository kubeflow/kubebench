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
	"path"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kubeflow/kubebench/controller/pkg/util"
)

const (
	experimentIDEnvName = "KUBEBENCH_EXP_ID"
)

type Configurator struct {
	FileOperator      FileOperatorInterface
	ManifestGenerator ManifestGeneratorInterface
	ManifestModifier  ManifestModifierInterface
}

func (c *Configurator) Run(options *AppOption) error {

	config := options.Config
	namespace := options.Namespace
	manifestOutput := options.ManifestOutput
	experimentIDOutput := options.ExperimentIDOutput
	var ksPrototypeRef KsPrototypeRef
	if err := json.Unmarshal([]byte(options.TemplateRef), &ksPrototypeRef); err != nil {
		log.Errorf("Cannot unmarshal value: %s", options.TemplateRef)
		return err
	}
	var ownerReferences []metav1.OwnerReference
	if err := json.Unmarshal([]byte(options.OwnerReferences), &ownerReferences); err != nil {
		log.Errorf("Cannot unmarshal value: %s", options.OwnerReferences)
		return err
	}
	var volumes []corev1.Volume
	if err := json.Unmarshal([]byte(options.Volumes), &volumes); err != nil {
		log.Errorf("Cannot unmarshal value: %s", options.Volumes)
		return err
	}
	var volumeMounts []corev1.VolumeMount
	if err := json.Unmarshal([]byte(options.VolumeMounts), &volumeMounts); err != nil {
		log.Errorf("Cannot unmarshal value: %s", options.VolumeMounts)
		return err
	}
	var envVars []corev1.EnvVar
	if err := json.Unmarshal([]byte(options.EnvVars), &envVars); err != nil {
		log.Errorf("Cannot unmarshal value: %s", options.EnvVars)
		return err
	}

	// Read and parse config
	parametersRaw, err := c.FileOperator.ReadConfig(config)
	if err != nil {
		log.Errorf("Failed to read runner config: %s", err)
		return err
	}
	var parameters map[string]interface{}
	if err := yaml.Unmarshal(parametersRaw, &parameters); err != nil {
		log.Errorf("Could not parse job parameters; Error: %s\n", err)
		return err
	}

	// Generate experiment ID
	experimentName := strings.TrimSuffix(config, path.Ext(config))
	experimentID := experimentName + "-" + time.Now().Format("200601021504") + "-" + util.RandString(4)
	// modify environment variable with experiment id
	for i, env := range envVars {
		if env.Name == experimentIDEnvName {
			envVars[i].Value = experimentID
		}
	}

	// Generate manifest
	manifest, err := c.ManifestGenerator.GenerateManifest(ksPrototypeRef, parameters)
	if err != nil {
		log.Errorf("Failed to generate manifest: %s", err)
		return err
	}

	// Modify manifest
	modSpec := ManifestModSpec{
		Name:            experimentID,
		Namespace:       namespace,
		OwnerReferences: ownerReferences,
		Volumes:         volumes,
		VolumeMounts:    volumeMounts,
		EnvVars:         envVars,
	}
	manifest, err = c.ManifestModifier.ModifyManifest(manifest, modSpec)
	if err != nil {
		log.Errorf("Failed to modify manifest: %s", err)
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
		configFilename:   parametersRaw,
		manifestFilename: manifest,
	}
	err = c.FileOperator.InitExperiment(experimentID, outputsMap)
	if err != nil {
		log.Errorf("Failed to initialize experiment: %s", err)
		return err
	}

	return nil
}
