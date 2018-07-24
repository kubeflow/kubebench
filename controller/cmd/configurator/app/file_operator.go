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
	"io/ioutil"
	"os"
	"path"

	"github.com/ghodss/yaml"
	log "github.com/sirupsen/logrus"

	"github.com/kubeflow/kubebench/controller/pkg/apis/kubebench/v1alpha1"
	"github.com/kubeflow/kubebench/controller/pkg/util"
)

type FileOperatorInterface interface {
	ReadRunnerConfig(runnerConfigFile string) (*v1alpha1.RunnerConfig, error)
	WriteManifest(manifest []byte, manifestFile string) error
}

type FileOperator struct{}

func (fo *FileOperator) ReadRunnerConfig(runnerConfigFile string) (*v1alpha1.RunnerConfig, error) {
	runnerConfig := &v1alpha1.RunnerConfig{}
	log.Printf("Loading job runner config from %s.", runnerConfigFile)
	data, err := ioutil.ReadFile(runnerConfigFile)
	if err != nil {
		log.Errorf("Could not read file: %s. Error: %s", runnerConfigFile, err)
		return nil, err
	}
	err = yaml.Unmarshal(data, runnerConfig)
	if err != nil {
		log.Errorf("Could not parse job runner config; Error: %s\n", err)
		return nil, err
	}
	log.Infof("RunnerConfig: %s", util.Pformat(runnerConfig))

	return runnerConfig, nil
}

func (fo *FileOperator) WriteManifest(manifest []byte, manifestFile string) error {
	dir, _ := path.Split(manifestFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Errorf("Failed to create directory: %s", err)
		return err
	}
	if err := ioutil.WriteFile(manifestFile, manifest, 0644); err != nil {
		log.Errorf("Failed to write output file: %s", err)
		return err
	}
	return nil
}
