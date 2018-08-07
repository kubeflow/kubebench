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
	"path"

	log "github.com/sirupsen/logrus"

	"github.com/kubeflow/kubebench/controller/pkg/apis/kubebench/v1alpha1"
	"github.com/kubeflow/kubebench/controller/pkg/util"
)

type ManifestGeneratorInterface interface {
	GenerateManifest(runnerConfig *v1alpha1.RunnerConfig) ([]byte, error)
}

type ManifestGenerator struct{}

func (mg *ManifestGenerator) ksInit(name string, cwd string) error {
	_, err := util.Run([]string{"ks", "init", name}, cwd, nil)
	if err != nil {
		log.Errorf("Failed to initialize Ksonnet app: %v", err)
		return err
	}
	return nil
}

func (mg *ManifestGenerator) ksRegistryAdd(name string, path string, cwd string) error {
	_, err := util.Run([]string{"ks", "registry", "add", name, path}, cwd, nil)
	if err != nil {
		log.Errorf("Failed to add Ksonnet registry: %v", err)
		return err
	}
	return nil
}

func (mg *ManifestGenerator) ksPkgInstall(name string, cwd string) error {
	_, err := util.Run([]string{"ks", "pkg", "install", name}, cwd, nil)
	if err != nil {
		log.Errorf("Failed to install Ksonnet package: %v", err)
		return err
	}
	return nil
}

func (mg *ManifestGenerator) ksGenerate(name string, cwd string) error {
	_, err := util.Run([]string{"ks", "generate", name, name}, cwd, nil)
	if err != nil {
		log.Errorf("Failed to generate Ksonnet component: %v", err)
		return err
	}
	return nil
}

func (mg *ManifestGenerator) ksParamsSet(name string, params map[string]interface{}, cwd string) error {
	for key, val := range params {
		if val == nil {
			continue
		}
		_, err := util.Run(
			[]string{"ks", "param", "set", name, key, "--", fmt.Sprintf("%v", val)}, cwd, nil)
		if err != nil {
			log.Errorf("Failed to set Ksonnet parameter: %v", err)
			return err
		}
	}
	return nil
}

func (mg *ManifestGenerator) ksShow(name string, cwd string) ([]byte, error) {
	out, err := util.Run([]string{"ks", "show", "default", "-c", name}, cwd, nil)
	if err != nil {
		log.Errorf("Failed to show manifest: %v", err)
		return nil, err
	}
	return out, nil
}

func (mg *ManifestGenerator) GenerateManifest(runnerConfig *v1alpha1.RunnerConfig) ([]byte, error) {
	wkdir := "/tmp"
	ksApp := "kubebench-app"
	ksdir := path.Join(wkdir, ksApp)

	ksPrototype := runnerConfig.Spec.Prototype.Name
	ksPackage := runnerConfig.Spec.Prototype.Package
	ksRegistryPath := runnerConfig.Spec.Prototype.Registry
	_, ksRegistry := path.Split(path.Clean(ksRegistryPath))
	params := runnerConfig.Spec.Parameters

	if err := mg.ksInit(ksApp, wkdir); err != nil {
		return nil, err
	}
	if err := mg.ksRegistryAdd(ksRegistry, ksRegistryPath, ksdir); err != nil {
		return nil, err
	}
	if err := mg.ksPkgInstall(ksRegistry+"/"+ksPackage, ksdir); err != nil {
		return nil, err
	}
	if err := mg.ksGenerate(ksPrototype, ksdir); err != nil {
		return nil, err
	}
	if err := mg.ksParamsSet(ksPrototype, params, ksdir); err != nil {
		return nil, err
	}
	manifest, err := mg.ksShow(ksPrototype, ksdir)
	if err != nil {
		return nil, err
	}

	return manifest, nil
}
