// Copyright 2019 The Kubeflow Authors
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

package workflow

import (
	"fmt"
	"time"

	argov1alpha1 "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	kbjobv1alpha2 "github.com/kubeflow/kubebench/controller/pkg/apis/kubebenchjob/v1alpha2"
	"github.com/kubeflow/kubebench/controller/pkg/util"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kubeflow/kubebench/controller/pkg/constants"
)

type workflowInfo struct {
	experimentID        string
	namespace           string
	labels              map[string]string
	ownerReferences     []metav1.OwnerReference
	env                 []corev1.EnvVar
	volumes             []corev1.Volume
	volumeMap           map[string]corev1.Volume
	managedVolumes      []corev1.Volume
	managedVolumeMounts []corev1.VolumeMount
	managedVolumeMap    map[string]corev1.Volume
}

func newWorkflowInfo(kbjob *kbjobv1alpha2.KubebenchJob) *workflowInfo {

	// Create an easy-to-read unique experiment ID for each run of the workflow
	// The experiment ID will be placed in the label of all resources created by the workflow
	experimentID := kbjob.Name + "-" + time.Now().Format("0601021504") + "-" + util.RandString(4)

	ownerRefs := []metav1.OwnerReference{
		{
			APIVersion: argov1alpha1.SchemeGroupVersion.Group + "/" + argov1alpha1.SchemeGroupVersion.Version,
			Kind:       "Workflow",
			Name:       "{{workflow.name}}",
			UID:        "{{workflow.uid}}",
		},
	}

	labels := map[string]string{
		"kubebench.kubeflow.org/experiment-id": experimentID,
	}

	envVars := []corev1.EnvVar{
		{
			Name:  constants.ExpIDEnvName,
			Value: experimentID,
		},
		// NOTE: WorkflowRootPath is not available to use. Mount point at WorkflowExpRootPath.
		// {
		// 	Name:  constants.WorkflowRootEnvName,
		// 	Value: constants.WorkflowRootPath,
		// },
		{
			Name:  constants.WorkflowExpRootEnvName,
			Value: constants.WorkflowExpRootPath,
		},
		{
			Name:  constants.WorkflowExpPathEnvName,
			Value: fmt.Sprintf(constants.WorkflowExpPathFmt, experimentID),
		},
		{
			Name:  constants.ExpRootEnvName,
			Value: constants.ExpRootPath,
		},
		{
			Name:  constants.ExpPathEnvName,
			Value: fmt.Sprintf(constants.ExpPathFmt, experimentID),
		},
		{
			Name:  constants.ExpConfigPathEnvName,
			Value: fmt.Sprintf(constants.ExpConfigPathFmt, experimentID),
		},
		{
			Name:  constants.ExpOutputPathEnvName,
			Value: fmt.Sprintf(constants.ExpOutputPathFmt, experimentID),
		},
		{
			Name:  constants.ExpResultPathEnvName,
			Value: fmt.Sprintf(constants.ExpResultPathFmt, experimentID),
		},
	}

	volMap := map[string]corev1.Volume{}
	for _, v := range kbjob.Spec.Volumes {
		volMap[v.Name] = v
	}
	managedVols := []corev1.Volume{}
	managedVolMnts := []corev1.VolumeMount{}
	managedVolMap := map[string]corev1.Volume{}
	managedVolCands := []*corev1.Volume{
		kbjob.Spec.ManagedVolumes.ExperimentVolume,
		kbjob.Spec.ManagedVolumes.WorkflowVolume,
	}
	managedVolMntPaths := []string{
		constants.ExpRootPath,
		constants.WorkflowExpRootPath,
	}
	for i, v := range managedVolCands {
		if v != nil {
			managedVols = append(managedVols, *v)
			volMnt := corev1.VolumeMount{
				Name:      v.Name,
				MountPath: managedVolMntPaths[i],
			}
			managedVolMnts = append(managedVolMnts, volMnt)
			managedVolMap[v.Name] = *v
		}
	}

	wfInfo := &workflowInfo{
		experimentID:        experimentID,
		namespace:           kbjob.Namespace,
		ownerReferences:     ownerRefs,
		labels:              labels,
		env:                 envVars,
		volumes:             kbjob.Spec.Volumes,
		volumeMap:           volMap,
		managedVolumes:      managedVols,
		managedVolumeMounts: managedVolMnts,
		managedVolumeMap:    managedVolMap,
	}

	return wfInfo
}
