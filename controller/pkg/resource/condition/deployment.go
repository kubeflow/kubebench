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

	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// DeploymentV1beta1Condition is a condition checker for deployments in extensions/v1beta1
type DeploymentV1beta1Condition struct{}

// NewDeploymentV1beta1Condition creates a new DeploymentV1beta1Condition
func NewDeploymentV1beta1Condition() *DeploymentV1beta1Condition {
	return &DeploymentV1beta1Condition{}
}

// CheckCondition checks the status of a given deployment.
// The success condition is met when number of available replicas = number of required replicas.
// The failure condition is met when there is any replica failure.
func (c *DeploymentV1beta1Condition) CheckCondition(resource *unstructured.Unstructured) (ResourceConditionStatus, error) {
	var deployment extensionsv1beta1.Deployment
	resStr, err := json.Marshal(resource)
	if err != nil {
		return ResourceConditionUnknown, err
	}
	err = json.Unmarshal(resStr, &deployment)
	if err != nil {
		return ResourceConditionUnknown, err
	}

	var result ResourceConditionStatus
	for _, cond := range deployment.Status.Conditions {
		if cond.Type == extensionsv1beta1.DeploymentReplicaFailure && cond.Status == corev1.ConditionTrue {
			result = ResourceConditionFailure
			break
		}
	}
	if deployment.Status.AvailableReplicas == deployment.Status.Replicas {
		result = ResourceConditionSuccess
	}

	return result, nil
}
