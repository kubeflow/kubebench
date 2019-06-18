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

package mod

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/kubeflow/kubebench/controller/pkg/resource/common"
)

var newResourceModFuncs = map[string]func() ResourceModifierInterface{
	"Job.v1.batch":                  NewJobV1Modifier,
	"MPIJob.v1alpha2":               NewMPIJobV1alpha2Modifier,
	"Deployment.v1beta1.extensions": NewDeploymentV1beta1Modifier,
}

type ResourceModifierInterface interface {
	ModifyResource(res *unstructured.Unstructured, modSpec *ResourceModSpec) (*unstructured.Unstructured, error)
}

type ResourceModifier struct{}

func NewResourceModifier() ResourceModifierInterface {
	modifier := &ResourceModifier{}
	return modifier
}

func (m *ResourceModifier) ModifyResource(
	res *unstructured.Unstructured,
	modSpec *ResourceModSpec) (*unstructured.Unstructured, error) {

	newRes := res.DeepCopy()

	// apply metadata modification to all resources
	newRes, err := NewMetaV1Modifier().ModifyResource(newRes, modSpec)
	// apply Kind-specific changes to supported kinds and versions
	// if Kind is not in the supported list, return without further change
	kvg := common.NewResourceRefFromUnstructured(res).SprintKindVersionGroup()
	newModFunc, found := newResourceModFuncs[kvg]
	if !found {
		return newRes, nil
	}
	modifier := newModFunc()
	newRes, err = modifier.ModifyResource(newRes, modSpec)
	if err != nil {
		return nil, err
	}

	return newRes, nil
}
