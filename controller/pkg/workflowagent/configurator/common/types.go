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

package common

import (
	"github.com/kubeflow/kubebench/controller/pkg/resource/mod"

	"github.com/kubeflow/kubebench/controller/pkg/apis/kubebenchjob/v1alpha2"
)

// ConfiguratorInput is the specification of configurator input
type ConfiguratorInput struct {
	*ManifestGenSpec `json:",inline"`
	*ManifestModSpec `json:",inline"`
}

// ManifestGenSpec is the spec for a manifest generation
type ManifestGenSpec struct {
	Manifest     *string                  `json:"manifest,omitempty"`
	ManifestFrom *v1alpha2.ManifestSource `json:"manifestFrom,omitempty"`
}

// ManifestModSpec is the spec of a manifest modification
type ManifestModSpec mod.ResourceModSpec
