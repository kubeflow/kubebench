// Copyright 2018 The Kubeflow Authors
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

package testutil

import (
	"time"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	common "github.com/kubeflow/tf-operator/pkg/apis/common/v1beta1"
	tfv1beta1 "github.com/kubeflow/tf-operator/pkg/apis/tensorflow/v1beta1"
)

func NewTFJobWithCleanPolicy(chief, worker, ps int, policy common.CleanPodPolicy) *tfv1beta1.TFJob {
	if chief == 1 {
		tfJob := NewTFJobWithChief(worker, ps)
		tfJob.Spec.CleanPodPolicy = &policy
		return tfJob
	}
	tfJob := NewTFJob(worker, ps)
	tfJob.Spec.CleanPodPolicy = &policy
	return tfJob
}

func NewTFJobWithCleanupJobDelay(chief, worker, ps int, ttl *int32) *tfv1beta1.TFJob {
	if chief == 1 {
		tfJob := NewTFJobWithChief(worker, ps)
		tfJob.Spec.TTLSecondsAfterFinished = ttl
		policy := common.CleanPodPolicyNone
		tfJob.Spec.CleanPodPolicy = &policy
		return tfJob
	}
	tfJob := NewTFJob(worker, ps)
	tfJob.Spec.TTLSecondsAfterFinished = ttl
	policy := common.CleanPodPolicyNone
	tfJob.Spec.CleanPodPolicy = &policy
	return tfJob
}

func NewTFJobWithChief(worker, ps int) *tfv1beta1.TFJob {
	tfJob := NewTFJob(worker, ps)
	tfJob.Spec.TFReplicaSpecs[tfv1beta1.TFReplicaTypeChief] = &common.ReplicaSpec{
		Template: NewTFReplicaSpecTemplate(),
	}
	return tfJob
}

func NewTFJobWithEvaluator(worker, ps, evaluator int) *tfv1beta1.TFJob {
	tfJob := NewTFJob(worker, ps)
	if evaluator > 0 {
		evaluator := int32(evaluator)
		tfJob.Spec.TFReplicaSpecs[tfv1beta1.TFReplicaTypeEval] = &common.ReplicaSpec{
			Replicas: &evaluator,
			Template: NewTFReplicaSpecTemplate(),
		}
	}
	return tfJob
}

func NewTFJob(worker, ps int) *tfv1beta1.TFJob {
	tfJob := &tfv1beta1.TFJob{
		TypeMeta: metav1.TypeMeta{
			Kind: tfv1beta1.Kind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      TestTFJobName,
			Namespace: metav1.NamespaceDefault,
		},
		Spec: tfv1beta1.TFJobSpec{
			TFReplicaSpecs: make(map[tfv1beta1.TFReplicaType]*common.ReplicaSpec),
		},
	}

	if worker > 0 {
		worker := int32(worker)
		workerReplicaSpec := &common.ReplicaSpec{
			Replicas: &worker,
			Template: NewTFReplicaSpecTemplate(),
		}
		tfJob.Spec.TFReplicaSpecs[tfv1beta1.TFReplicaTypeWorker] = workerReplicaSpec
	}

	if ps > 0 {
		ps := int32(ps)
		psReplicaSpec := &common.ReplicaSpec{
			Replicas: &ps,
			Template: NewTFReplicaSpecTemplate(),
		}
		tfJob.Spec.TFReplicaSpecs[tfv1beta1.TFReplicaTypePS] = psReplicaSpec
	}
	return tfJob
}

func NewTFReplicaSpecTemplate() v1.PodTemplateSpec {
	return v1.PodTemplateSpec{
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				v1.Container{
					Name:  tfv1beta1.DefaultContainerName,
					Image: TestImageName,
					Args:  []string{"Fake", "Fake"},
					Ports: []v1.ContainerPort{
						v1.ContainerPort{
							Name:          tfv1beta1.DefaultPortName,
							ContainerPort: tfv1beta1.DefaultPort,
						},
					},
				},
			},
		},
	}
}

func SetTFJobCompletionTime(tfJob *tfv1beta1.TFJob) {
	now := metav1.Time{Time: time.Now()}
	tfJob.Status.CompletionTime = &now
}
