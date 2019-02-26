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
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/kubeflow/kubebench/controller/pkg/resource/common"
)

var supportedResourceConditions = map[string]ResourceConditionInterface{
	"Job.v1.batch":                  NewJobV1Condition(),
	"Deployment.v1beta1.extensions": NewDeploymentV1beta1Condition(),
}

// ResourceConditionStatus is the status of a resource condition check
type ResourceConditionStatus string

const (
	// ResourceConditionSuccess indicates success condition met
	ResourceConditionSuccess ResourceConditionStatus = "Success"
	// ResourceConditionFailure indicates failure condition met
	ResourceConditionFailure ResourceConditionStatus = "Failure"
	// ResourceConditionTimeout indicates condition check timed out
	ResourceConditionTimeout ResourceConditionStatus = "Timeout"
	// ResourceConditionUnknown indicates condition check returned unknown status or error
	ResourceConditionUnknown ResourceConditionStatus = "Unknown"
)

// ResourceConditionInterface is the interface of resource condition checkers
type ResourceConditionInterface interface {
	CheckCondition(resource *unstructured.Unstructured) (ResourceConditionStatus, error)
}

// NewResourceCondition creates a new resource condition checker
// The type of the condition checker depends on the given resource
func NewResourceCondition(ref *common.ResourceRef) ResourceConditionInterface {
	kvg := ref.SprintKindVersionGroup()
	rc, found := supportedResourceConditions[kvg]
	if found == false {
		rc = nil
	}
	return rc
}
