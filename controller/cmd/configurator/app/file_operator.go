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
	"fmt"
	"io/ioutil"
	"os"
	"path"

	log "github.com/sirupsen/logrus"

	"github.com/kubeflow/kubebench/controller/pkg/constants"
)

type FileOperatorInterface interface {
	ReadConfig(configFile string) ([]byte, error)
	WriteOutputs(outputsMap map[string][]byte) error
	InitExperiment(experimentID string, outputsMap map[string][]byte) error
}

type FileOperator struct{}

func (fo *FileOperator) ReadConfig(configFile string) ([]byte, error) {
	configRoot := os.Getenv("KUBEBENCH_CONFIG_ROOT")
	data, err := ioutil.ReadFile(path.Join(configRoot, configFile))
	if err != nil {
		log.Errorf("Could not read file: %s. Error: %s", configFile, err)
		return nil, err
	}
	return data, nil
}

func (fo *FileOperator) WriteOutputs(outputsMap map[string][]byte) error {
	for file, data := range outputsMap {
		if err := fo.writeFileNewDir(data, file); err != nil {
			return err
		}
	}
	return nil
}

func (fo *FileOperator) InitExperiment(experimentID string, outputsMap map[string][]byte) error {
	expID := os.Getenv(constants.ExpIDEnvName)
	expConfigDir := fmt.Sprintf(constants.ExpConfigPathFmt, expID)
	expOutputDir := fmt.Sprintf(constants.ExpOutputPathFmt, expID)
	expResultDir := fmt.Sprintf(constants.ExpResultPathFmt, expID)
	wfExpDir := fmt.Sprintf(constants.WorkflowExpPathFmt, expID)
	dirsToCreate := [...]string{expConfigDir, expOutputDir, expResultDir, wfExpDir}
	for _, dir := range dirsToCreate {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	for file, data := range outputsMap {
		if err := fo.writeFileNewDir(data, path.Join(expConfigDir, file)); err != nil {
			return err
		}
	}

	return nil
}

func (fo *FileOperator) writeFileNewDir(data []byte, file string) error {
	dir, _ := path.Split(file)
	if dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Errorf("Failed to create directory: %s. Error: %s", dir, err)
			return err
		}
	}
	if err := ioutil.WriteFile(file, data, 0644); err != nil {
		log.Errorf("Failed to write output file: %s. Error: %s", file, err)
		return err
	}
	return nil
}
