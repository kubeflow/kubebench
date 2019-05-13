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
	corev1 "k8s.io/api/core/v1"
)

// ModifyPodTemplateV1 modifies a pod template with given changes
func ModifyPodTemplateV1(template corev1.PodTemplateSpec, modSpec *ResourceModSpec) corev1.PodTemplateSpec {
	newTemplate := template
	newTemplate.Spec.Volumes = append(template.Spec.Volumes, modSpec.Volumes...)
	for i, container := range newTemplate.Spec.Containers {
		newContainer := ModifyContainerV1(container, modSpec)
		newTemplate.Spec.Containers[i] = newContainer
	}
	return newTemplate
}

// ModifyContainerV1 modifies a container with given changes
func ModifyContainerV1(container corev1.Container, modSpec *ResourceModSpec) corev1.Container {
	newContainer := container
	newContainer.VolumeMounts = append(container.VolumeMounts, modSpec.VolumeMounts...)
	newContainer.Env = append(container.Env, modSpec.Env...)
	return newContainer
}
