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
	argov1alpha1 "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/imdario/mergo"
	kbjobv1alpha2 "github.com/kubeflow/kubebench/controller/pkg/apis/kubebenchjob/v1alpha2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// BuildWorkflow builds an Argo Workflow from a KubebenchJob
func BuildWorkflow(
	kbjobIn *kbjobv1alpha2.KubebenchJob,
	kbconfig *kbjobv1alpha2.KubebenchConfig,
	inOperator bool) (*argov1alpha1.Workflow, error) {

	kbjob := kbjobIn.DeepCopy()

	// Merge the KubebenchJob with default Kubebench config
	if err := applyKubebenchConfig(kbjob, kbconfig); err != nil {
		return nil, err
	}

	wfInfo := newWorkflowInfo(kbjob)

	// If in Kubebench operator, then set workflow owner to the KubebenchJob
	var ownerRefs []metav1.OwnerReference
	if inOperator {
		ownerRefs = append(
			ownerRefs,
			*metav1.NewControllerRef(kbjob, schema.GroupVersionKind{
				Group:   kbjobv1alpha2.GroupName,
				Version: kbjobv1alpha2.GroupVersion,
				Kind:    kbjobv1alpha2.Kind,
			}))
	}

	metadata := metav1.ObjectMeta{
		Name:            kbjob.Name,
		Namespace:       wfInfo.namespace,
		OwnerReferences: ownerRefs,
		Labels:          wfInfo.labels,
	}

	var workflowTemplates []argov1alpha1.Template
	var dagTasks []argov1alpha1.DAGTask
	depMap := map[string]string{}
	tasks := kbjob.Spec.Tasks
	for _, task := range tasks {
		if task.Container != nil {
			wfTemplate := buildContainerTemplate(
				task.Name, task.Container, wfInfo, task.Inputs, task.Outputs)
			workflowTemplates = append(workflowTemplates, wfTemplate)
			wfTask := buildDAGTask(
				wfTemplate.Name, task.Dependencies, kbjob.Spec.Arguments,
				task.Inputs, task.Outputs)
			dagTasks = append(dagTasks, wfTask)
			depMap[task.Name] = wfTask.Name
		} else if task.Resource != nil {
			// For resource typed tasks, the inputs are passed to configurator.
			// This makes it easy to use inputs for configuration sources.
			configTemplate := buildResourceConfigTemplate(
				task.Name, kbjob.Spec.WorkflowAgent.Container, task.Resource,
				wfInfo, task.Inputs)
			workflowTemplates = append(workflowTemplates, configTemplate)
			configTask := buildDAGTask(
				configTemplate.Name, task.Dependencies, kbjob.Spec.Arguments,
				task.Inputs, argov1alpha1.Outputs{})
			dagTasks = append(dagTasks, configTask)
			lastTask := configTask.Name

			createTemplate := buildResourceCreateTemplate(
				task.Name, kbjob.Spec.WorkflowAgent.Container, task.Resource, wfInfo)
			workflowTemplates = append(workflowTemplates, createTemplate)
			createTask := buildDAGTask(
				createTemplate.Name, []string{configTask.Name}, kbjob.Spec.Arguments,
				argov1alpha1.Inputs{}, argov1alpha1.Outputs{})
			dagTasks = append(dagTasks, createTask)
			lastTask = createTask.Name

			if task.Resource.Options != nil && task.Resource.Options.AutoWatch != nil {
				autoWatchTemplate := buildResourceAutoWatchTemplate(
					task.Name, kbjob.Spec.WorkflowAgent.Container, task.Resource, wfInfo)
				workflowTemplates = append(workflowTemplates, autoWatchTemplate)
				autoWatchTask := buildDAGTask(
					autoWatchTemplate.Name, []string{createTask.Name}, kbjob.Spec.Arguments,
					argov1alpha1.Inputs{}, argov1alpha1.Outputs{})
				dagTasks = append(dagTasks, autoWatchTask)
				lastTask = autoWatchTask.Name
			}

			// Create a separate DAG task to handle outputs
			if task.Outputs.Parameters != nil || task.Outputs.Artifacts != nil {
				wfTemplate := buildIOTemplate(
					task.Name, kbjob.Spec.WorkflowAgent.Container, wfInfo,
					argov1alpha1.Inputs{}, task.Outputs)
				workflowTemplates = append(workflowTemplates, wfTemplate)
				outputTask := buildDAGTask(
					wfTemplate.Name, []string{lastTask}, kbjob.Spec.Arguments,
					argov1alpha1.Inputs{}, task.Outputs)
				dagTasks = append(dagTasks, outputTask)
				lastTask = outputTask.Name
			}

			depMap[task.Name] = lastTask
		}
	}
	// Map dependencies from kubebenchjob task to workflow task
	for i, task := range dagTasks {
		for j, dep := range task.Dependencies {
			if newDep, found := depMap[dep]; found {
				dagTasks[i].Dependencies[j] = newDep
			}
		}
	}

	dagTemplate := buildDAGTemplate(
		"kubebench-job-workflow-entrypoint", dagTasks, kbjob.Spec.Arguments)

	workflowTemplates = append(workflowTemplates, dagTemplate)

	workflow := &argov1alpha1.Workflow{
		TypeMeta: metav1.TypeMeta{
			APIVersion: argov1alpha1.SchemeGroupVersion.Group + "/" + argov1alpha1.SchemeGroupVersion.Version,
			Kind:       "Workflow",
		},
		ObjectMeta: metadata,
		Spec: argov1alpha1.WorkflowSpec{
			ServiceAccountName: kbjob.Spec.ServiceAccountName,
			Arguments:          kbjob.Spec.Arguments,
			Entrypoint:         "kubebench-job-workflow-entrypoint",
			Templates:          workflowTemplates,
			Volumes:            append(wfInfo.volumes, wfInfo.managedVolumes...),
		},
	}

	return workflow, nil
}

// applyKubebenchConfig merges the KubebenchJob with default values in KubebenchConfig,
// the existing fields in the KubebenchJob will take priority
func applyKubebenchConfig(
	kbjob *kbjobv1alpha2.KubebenchJob, kbconfig *kbjobv1alpha2.KubebenchConfig) error {

	// managed volumes are set to default if not set in kubebench job spec,
	// the fields inside volumes should not be merged recursively
	if kbjob.Spec.ManagedVolumes.ExperimentVolume == nil {
		kbjob.Spec.ManagedVolumes.ExperimentVolume = kbconfig.DefaultManagedVolumes.ExperimentVolume
	}
	if kbjob.Spec.ManagedVolumes.WorkflowVolume == nil {
		kbjob.Spec.ManagedVolumes.WorkflowVolume = kbconfig.DefaultManagedVolumes.WorkflowVolume
	}

	// workflow agent spec is merged recursively with default
	if err := mergo.Merge(&kbjob.Spec.WorkflowAgent, &kbconfig.DefaultWorkflowAgent); err != nil {
		return err
	}

	return nil
}
