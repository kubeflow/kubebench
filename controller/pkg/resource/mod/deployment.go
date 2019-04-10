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
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

type DeploymentV1beta1Modifier struct{}

func NewDeploymentV1beta1Modifier() ResourceModifierInterface {
	modifier := &DeploymentV1beta1Modifier{}
	return modifier
}

func (m *DeploymentV1beta1Modifier) ModifyResource(
	res *unstructured.Unstructured,
	modSpec *ResourceModSpec) (*unstructured.Unstructured, error) {

	deployment := &extensionsv1beta1.Deployment{}
	converter := runtime.DefaultUnstructuredConverter
	if err := converter.FromUnstructured(res.Object, deployment); err != nil {
		return nil, err
	}

	deployment.Spec.Template = ModifyPodTemplateV1(deployment.Spec.Template, modSpec)

	newResObj, err := converter.ToUnstructured(deployment)
	if err != nil {
		return nil, err
	}
	newRes := &unstructured.Unstructured{Object: newResObj}

	return newRes, nil
}
