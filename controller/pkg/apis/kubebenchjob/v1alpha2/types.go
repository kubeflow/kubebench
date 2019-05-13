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

package v1alpha2

import (
	argov1alpha1 "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +resource:path=kubebenchjob

// KubebenchJob is the definition of a KubebenchJob resource
type KubebenchJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KubebenchJobSpec   `json:"spec,omitempty"`
	Status KubebenchJobStatus `json:"status,omitempty"`
}

// KubebenchJobSpec is the specification of a KubebenchJob
type KubebenchJobSpec struct {
	ServiceAccountName string            `json:"serviceAccountName,omitempty"`
	Volumes            []corev1.Volume   `json:"volumes,omitempty"`
	ManagedVolumes     ManagedVolumes    `json:"managedVolumes,omitempty"`
	WorkflowAgent      WorkflowAgentSpec `json:"workflowAgent,omitempty"`
	Tasks              []Task            `json:"tasks"`
}

// ManagedVolumes is the volumes managed by the Kubebench workflow
type ManagedVolumes struct {
	ExperimentVolume *corev1.Volume `json:"experimentVolume,omitempty"`
	WorkflowVolume   *corev1.Volume `json:"workflowVolume,omitempty"`
}

// WorkflowAgentSpec is the specification of Kubebench workflow agent
type WorkflowAgentSpec struct {
	Container *corev1.Container `json:"container,omitempty"`
}

// Task is a task in a Kubebench workflow
type Task struct {
	Name         string            `json:"name"`
	Container    *corev1.Container `json:"container,omitempty"`
	Resource     *ResourceSpec     `json:"resource,omitempty"`
	Dependencies []string          `json:"dependencies,omitempty"`
}

// ResourceSpec is the specification of a resouce-type task in a benchmark workflow
type ResourceSpec struct {
	Manifest     *string              `json:"manifest,omitempty"`
	ManifestFrom *ManifestSource      `json:"manifestFrom,omitempty"`
	VolumeMounts []corev1.VolumeMount `json:"volumeMounts,omitempty"`
	Options      *ResourceOptions     `json:"options,omitempty"`
}

// ManifestSource is the source of a k8s manifest
type ManifestSource struct {
	Path *string `json:"path,omitempty"`
}

// ResourceOptions is the options of a resource-type task
type ResourceOptions struct {
	MountManagedVolumes bool           `json:"mountManagedVolumes,omitempty"`
	AutoDelete          bool           `json:"autoDelete,omitempty"`
	AutoWatch           *AutoWatchSpec `json:"autoWatch,omitempty"`
	NumCopies           int            `json:"numCopies,omitempty"`
}

// AutoWatchSpec is the specification of auto-watch functionality
type AutoWatchSpec struct {
	Timeout string `json:"timeout,omitempty"`
}

// KubebenchJobStatus is the observed status of a KubebenchJob
type KubebenchJobStatus struct {
	argov1alpha1.WorkflowStatus `json:",inline"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KubebenchJobList contains a list of KubebenchJob
type KubebenchJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KubebenchJob `json:"items"`
}

// KubebenchConfig is the Kubebench configuration
type KubebenchConfig struct {
	DefaultWorkflowAgent  WorkflowAgentSpec `json:"defaultWorkflowAgent,omitempty"`
	DefaultManagedVolumes ManagedVolumes    `json:"defaultManagedVolumes,omitempty"`
}
