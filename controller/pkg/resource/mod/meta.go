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

import "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

type MetaV1Modifier struct{}

func NewMetaV1Modifier() *MetaV1Modifier {
	modifier := &MetaV1Modifier{}
	return modifier
}

func (m *MetaV1Modifier) ModifyResource(
	res *unstructured.Unstructured,
	modSpec *ResourceModSpec) (*unstructured.Unstructured, error) {

	newRes := res.DeepCopy()
	// set namespace
	newRes.SetNamespace(modSpec.Namespace)
	// set owner references
	newOwnerReferences := append(res.GetOwnerReferences(), modSpec.OwnerReferences...)
	newRes.SetOwnerReferences(newOwnerReferences)
	// set labels
	newLabels := res.GetLabels()
	if newLabels == nil {
		newLabels = make(map[string]string)
	}
	for k, v := range modSpec.Labels {
		newLabels[k] = v
	}
	newRes.SetLabels(newLabels)

	return newRes, nil
}
