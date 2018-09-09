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
	batchv1 "k8s.io/api/batch/v1"
)

type JobManifestModifier struct {
	BaseManifestModifier
}

func (mm *JobManifestModifier) ModifyManifest(manifest []byte, modSpec ManifestModSpec) ([]byte, error) {
	manifest, err := mm.BaseManifestModifier.ModifyManifest(manifest, modSpec)
	if err != nil {
		return nil, err
	}
	var job batchv1.Job
	if err := yaml.Unmarshal(manifest, &job); err != nil {
		log.Errorf("Failed to unmarshal manifest: %s", manifest)
		return nil, err
	}
	job.Spec.Template.Spec.Volumes = append(job.Spec.Template.Spec.Volumes, modSpec.Volumes...)
	for i, container := range job.Spec.Template.Spec.Containers {
		job.Spec.Template.Spec.Containers[i].VolumeMounts = append(container.VolumeMounts, modSpec.VolumeMounts...)
		job.Spec.Template.Spec.Containers[i].Env = append(container.Env, modSpec.EnvVars...)
	}
	manifest, err = yaml.Marshal(job)
	if err != nil {
		log.Errorf("Failed to create modified manifest: %s", err)
		return nil, err
	}
	return manifest, nil
}
