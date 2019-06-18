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
	mpijob "github.com/kubeflow/mpi-operator/pkg/apis/kubeflow/v1alpha2"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

type JobV1Modifier struct{}

type MPIJobV1alpha2Modifier struct{}

func NewJobV1Modifier() ResourceModifierInterface {
	modifier := &JobV1Modifier{}
	return modifier
}

func NewMPIJobV1alpha2Modifier() ResourceModifierInterface {
	modifier := &MPIJobV1alpha2Modifier{}
	return modifier
}

func (m *JobV1Modifier) ModifyResource(
	res *unstructured.Unstructured,
	modSpec *ResourceModSpec) (*unstructured.Unstructured, error) {

	job := &batchv1.Job{}
	converter := runtime.DefaultUnstructuredConverter
	if err := converter.FromUnstructured(res.Object, job); err != nil {
		return nil, err
	}

	job.Spec.Template = ModifyPodTemplateV1(job.Spec.Template, modSpec)

	newResObj, err := converter.ToUnstructured(job)
	if err != nil {
		return nil, err
	}
	newRes := &unstructured.Unstructured{Object: newResObj}

	return newRes, nil
}

func (m *MPIJobV1alpha2Modifier) ModifyResource(
	res *unstructured.Unstructured,
	modSpec *ResourceModSpec) (*unstructured.Unstructured, error) {

	job := &mpijob.MPIJob{}
	newRes := &unstructured.Unstructured{}
	return newRes, nil

}
