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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type BaseJob struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata
	Metadata metav1.ObjectMeta `json:"metadata,omitempty"`
	// Specification of the job
	Spec interface{} `json:"spec,omitempty"`
}

// ManifestModSpec provides specs that the manifest should be modified with
type ManifestModSpec struct {
	// Resource metadata.name
	Name string
	// Resource metadata.namespace
	Namespace string
	// Resource metadata.ownerReferences
	OwnerReferences []metav1.OwnerReference
	// Pod volumes
	Volumes []corev1.Volume
	// Container volume mounts
	VolumeMounts []corev1.VolumeMount
	// Container environment variables
	EnvVars []corev1.EnvVar
}
