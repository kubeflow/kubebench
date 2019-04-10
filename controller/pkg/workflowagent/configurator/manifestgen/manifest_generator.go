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

package manifestgen

import (
	"errors"
	"reflect"

	"github.com/kubeflow/kubebench/controller/pkg/apis/kubebenchjob/v1alpha2"
	"github.com/kubeflow/kubebench/controller/pkg/workflowagent/configurator/common"
)

var newManifestGenFuncs = map[string]func(interface{}) ManifestGeneratorInterface{
	"Path": NewPathManifestGenerator,
}

// ManifestGeneratorInterface is the interface of a ManifestGenerator
type ManifestGeneratorInterface interface {
	GenerateManifest() ([]byte, error)
}

// ManifestGenerator generates a k8s manifest given a spec of manifest source
type ManifestGenerator struct {
	spec              *common.ManifestGenSpec
	newGeneratorFuncs map[string]func(interface{}) ManifestGeneratorInterface
}

// NewManifestGenerator creates a new ManifestGenerator
func NewManifestGenerator(spec *common.ManifestGenSpec) ManifestGeneratorInterface {
	generator := &ManifestGenerator{
		spec:              spec,
		newGeneratorFuncs: newManifestGenFuncs,
	}
	return generator
}

// GenerateManifest genretes a k8s manifest
func (mg *ManifestGenerator) GenerateManifest() ([]byte, error) {
	var manifest []byte
	var err error
	if mg.spec.Manifest != nil {
		manifest = []byte(*(mg.spec.Manifest))
	} else if mg.spec.ManifestFrom != nil {
		generator, err := mg.getManifestGenerator(mg.spec.ManifestFrom)
		manifest, err = generator.GenerateManifest()
		if err != nil {
			return nil, err
		}
	} else {
		err = errors.New("Invalid manifest generator spec")
		return nil, err
	}
	return manifest, nil
}

func (mg *ManifestGenerator) getManifestGenerator(source *v1alpha2.ManifestSource) (ManifestGeneratorInterface, error) {
	fields := reflect.TypeOf(*source)
	values := reflect.ValueOf(*source)
	numFields := 0
	var fieldName string
	var fieldValue interface{}
	for i := 0; i < fields.NumField(); i++ {
		if !values.Field(i).IsNil() {
			fieldName = fields.Field(i).Name
			fieldValue = values.Field(i).Interface()
			numFields++
		}
	}
	newFunc, newFuncExists := mg.newGeneratorFuncs[fieldName]
	if (numFields != 1) || (!newFuncExists) {
		err := errors.New("Invalid manifest source")
		return nil, err
	}
	manifestGenerator := newFunc(fieldValue)
	return manifestGenerator, nil
}
