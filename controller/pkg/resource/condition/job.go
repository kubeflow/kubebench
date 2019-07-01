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

package condition

import (
	"encoding/json"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	mpijob "github.com/kubeflow/mpi-operator/pkg/apis/kubeflow/v1alpha1"
)

// JobV1Condition is a condition checker for jobs in batch/v1
type JobV1Condition struct{}

// NewJobV1Condition creates a new JobV1Condition
func NewJobV1Condition() *JobV1Condition {
	return &JobV1Condition{}
}

type MPIJobV1alpha1Condition struct{}

func NewMPIJobV1alpha1Condition() *MPIJobV1alpha1Condition {
	return &MPIJobV1alpha1Condition{}
}

// CheckCondition checks the status of a given job.
// The success condition is met when a "Complete" type condition is observed.
// The failure condition is met when a "Failed" type condition is observed.
func (c *JobV1Condition) CheckCondition(resource *unstructured.Unstructured) (ResourceConditionStatus, error) {
	var job batchv1.Job
	resStr, err := json.Marshal(resource)
	if err != nil {
		return ResourceConditionUnknown, err
	}
	err = json.Unmarshal(resStr, &job)
	if err != nil {
		return ResourceConditionUnknown, err
	}

	var result ResourceConditionStatus
	for _, cond := range job.Status.Conditions {
		if cond.Type == batchv1.JobFailed && cond.Status == corev1.ConditionTrue {
			result = ResourceConditionFailure
			break
		} else if cond.Type == batchv1.JobComplete && cond.Status == corev1.ConditionTrue {
			result = ResourceConditionSuccess
			break
		}
	}

	return result, nil
}

func (c *MPIJobV1alpha1Condition) CheckCondition(resource *unstructured.Unstructured) (ResourceConditionStatus, error) {
	var job mpijob.MPIJob

	resStr, err := json.Marshal(resource)
	if err != nil {
		return ResourceConditionUnknown, err
	}

	err = json.Unmarshal(resStr, &job)
	if err != nil {
		return ResourceConditionUnknown, err
	}
	var result ResourceConditionStatus
	if job.Status.LauncherStatus == mpijob.LauncherFailed {
		result = ResourceConditionFailure
	} else if job.Status.LauncherStatus == mpijob.LauncherSucceeded {
		result = ResourceConditionSuccess
	}

	return result, nil
}
