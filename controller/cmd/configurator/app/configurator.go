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

	"github.com/ghodss/yaml"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kubeflow/kubebench/controller/pkg/apis/kubebench/v1alpha1"
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
	// modify environment variable with experiment id
	for i, env := range envVars {
		if env.Name == experimentIDEnvName {
			envVars[i].Value = experimentID
		}
	}

	// Generate manifest
	manifest, err := c.ManifestGenerator.GenerateManifest(runnerConfig)
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
