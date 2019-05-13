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

package v1alpha1

import (
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +resource:path=kubebenchjob

// KubebenchJob is the configuration of a Kubebench job
type KubebenchJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KubebenchJobSpec   `json:"spec,omitempty"`
	Status KubebenchJobStatus `json:"status,omitempty"`
}

// KubebenchJobSpec defines the desired state of KubebenchJob
type KubebenchJobSpec struct {
	ServiceAccount string      `json:"serviceAccount,omitempty"`
	Volumes        VolumeSpecs `json:"volumeSpecs"`
	Secrets        SecretSpecs `json:"secretSpecs,omitempty"`
	Jobs           JobSpecs    `json:"jobSpecs"`
	Report         ReportSpecs `json:"reportSpecs,omitempty"`
}

type VolumeSpecs struct {
	ConfigVolume     apiv1.Volume `json:"configVolume"`
	ExperimentVolume apiv1.Volume `json:"experimentVolume"`
	//add data volume
}

type SecretSpecs struct {
	GithubToken    map[string]string `json:"githubTokenSecret,omitempty"`
	GCPCredentials map[string]string `json:"gcpCredentialsSecret,omitempty"`
}

type JobSpecs struct {
	PreJob  JobSpec     `json:"preJob,omitempty"`
	MainJob MainJobSpec `json:"mainJob"`
	PostJob JobSpec     `json:"postJob,omitempty"`
}

type JobSpec struct {
	Container apiv1.Container `json:"container,omitempty"`
}

type MainJobSpec struct {
	Resource ResourceSpec `json:"resource"`
}

type ResourceSpec struct {
	ManifestTemplate   TemplateSpec   `json:"manifestTemplate"`
	MainfestParameters ParametersSpec `json:"manifestParameters"`
	CreateSuccess      string         `json:"createSuccessCondition"`
	CreateFailure      string         `json:"createFailureCondition"`
	RunSuccess         string         `json:"runSuccessCondition"`
	RunFailure         string         `json:"runFailureCondition"`
}

type TemplateSpec struct {
	ValueFrom Ksonnet `json:"valueFrom"`
}

type Ksonnet struct {
	KsonnetSpec map[string]string `json:"ksonnet"`
}

type ParametersSpec struct {
	ValueFrom PathSpec `json:"valueFrom"`
}

type PathSpec struct {
	Path string `json:"path"`
}

type ReportSpecs struct {
	CSV []map[string]string `json:"csv"`
}

// KubebenchJobStatus defines the observed state of KubebenchJob
type KubebenchJobStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KubebenchJobList contains a list of KubebenchJob
type KubebenchJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KubebenchJob `json:"items"`
}
