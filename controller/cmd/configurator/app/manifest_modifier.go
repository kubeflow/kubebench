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

	"github.com/ghodss/yaml"
	log "github.com/sirupsen/logrus"
)

type ManifestModifierInterface interface {
	ModifyManifest(manifest []byte, modSpec ManifestModSpec) ([]byte, error)
}

type ManifestModifierFactory struct {
	Name string
}

func (mmf *ManifestModifierFactory) NewManifestModifier() (ManifestModifierInterface, error) {
	var m ManifestModifierInterface
	switch mmf.Name {
	case "TFJob":
		m = &TFJobManifestModifier{}
	case "PyTorchJob":
		m = &PyTorchJobManifestModifier{}
	case "MPIJob":
		m = &MPIJobManifestModifier{}
	case "Job":
		m = &JobManifestModifier{}
	default:
		err := fmt.Errorf("Unknown manifest modifier: %s", mmf.Name)
		log.Errorf("%s", err)
		return nil, err
	}
	return m, nil
}

type ManifestModifier struct{}

func (mm *ManifestModifier) ModifyManifest(manifest []byte, modSpec ManifestModSpec) ([]byte, error) {
	var baseInfo struct {
		Kind string `json:"kind"`
	}
	if err := yaml.Unmarshal(manifest, &baseInfo); err != nil {
		return nil, err
	}
	kind := baseInfo.Kind
	manifestModifierFactory := ManifestModifierFactory{Name: kind}
	manifestModifier, err := manifestModifierFactory.NewManifestModifier()
	if err != nil {
		return nil, err
	}
	newManifest, err := manifestModifier.ModifyManifest(manifest, modSpec)
	if err != nil {
		return nil, err
	}
	return newManifest, nil
}
