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
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// ResourceRef provides information that identifies a resource in a cluster
type ResourceRef struct {
	Group     string `json:"group"`
	Version   string `json:"version"`
	Kind      string `json:"kind"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

// NewResourceRefFromUnstructured creates a new ResourceRef from an Unstructured object
func NewResourceRefFromUnstructured(resource *unstructured.Unstructured) *ResourceRef {
	ref := ResourceRef{
		Group:     resource.GroupVersionKind().Group,
		Kind:      resource.GroupVersionKind().Kind,
		Version:   resource.GroupVersionKind().Version,
		Name:      resource.GetName(),
		Namespace: resource.GetNamespace(),
	}
	return &ref
}

// SprintKindVersionGroup returns a string of "kind.version.group"
func (r *ResourceRef) SprintKindVersionGroup() string {
	return fmt.Sprintf("%s.%s.%s", r.Kind, r.Version, r.Group)
}

// SprintKindVersionGroupName returns a string of "kind.version.group/name"
func (r *ResourceRef) SprintKindVersionGroupName() string {
	return fmt.Sprintf("%s.%s.%s/%s", r.Kind, r.Version, r.Group, r.Name)
}

// SprintGroupVersion returns a string of "group/version"
func (r *ResourceRef) SprintGroupVersion() string {
	var groupVersion string
	if r.Group == "" {
		groupVersion = r.Version
	} else {
		groupVersion = r.Group + "/" + r.Version
	}
	return groupVersion
}
