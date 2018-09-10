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
	"github.com/ghodss/yaml"
	log "github.com/sirupsen/logrus"
)

type BaseManifestModifier struct{}

func (mm *BaseManifestModifier) ModifyManifest(manifest []byte, modSpec ManifestModSpec) ([]byte, error) {
	var baseJob BaseJob
	if err := yaml.Unmarshal(manifest, &baseJob); err != nil {
		log.Errorf("Failed to modify manifest: %s", err)
		return nil, err
	}
	baseJob.Metadata.Name = modSpec.Name
	baseJob.Metadata.Namespace = modSpec.Namespace
	baseJob.Metadata.OwnerReferences = modSpec.OwnerReferences
	manifest, err := yaml.Marshal(baseJob)
	if err != nil {
		log.Errorf("Failed to modify manifest: %s", err)
		return nil, err
	}
	return manifest, nil
}
