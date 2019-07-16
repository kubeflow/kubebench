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
	"encoding/json"
	"fmt"
	"path"
	"strconv"

	argov1alpha1 "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	kbjobv1alpha2 "github.com/kubeflow/kubebench/controller/pkg/apis/kubebenchjob/v1alpha2"
	corev1 "k8s.io/api/core/v1"

	"github.com/kubeflow/kubebench/controller/pkg/constants"
	"github.com/kubeflow/kubebench/controller/pkg/resource/mod"
	wfacommon "github.com/kubeflow/kubebench/controller/pkg/workflowagent/configurator/common"
)

const (
	configuratorTemplateNameFmt      = "%s-config"
	configuratorInputNameFmt         = configuratorTemplateNameFmt + "-in"
	configuratorOutputNameFmt        = configuratorTemplateNameFmt + "-out"
	resourceCreateTemplateNameFmt    = "%s-create"
	resourceCreateInputNameFmt       = resourceCreateTemplateNameFmt + "-in"
	resourceCreateOutputNameFmt      = resourceCreateTemplateNameFmt + "-out"
	resourceAutoWatchTemplateNameFmt = "%s-autowatch"
	resourceAutoWatchInputNameFmt    = resourceAutoWatchTemplateNameFmt + "-in"
	resourceAutoWatchOutputNameFmt   = resourceAutoWatchTemplateNameFmt + "-out"
)

func getName(format string, name string) string {
	return fmt.Sprintf(format, name)
}

func buildDAGTask(
	name string,
	dependencies []string,
	jobArguments argov1alpha1.Arguments,
	inputs argov1alpha1.Inputs,
	outputs argov1alpha1.Outputs) argov1alpha1.DAGTask {

	dagTask := argov1alpha1.DAGTask{
		Name:         name,
		Template:     name,
		Dependencies: dependencies,
	}

	argParameterNames := map[string]bool{}
	argArtifactNames := map[string]bool{}
	for _, p := range jobArguments.Parameters {
		argParameterNames[p.Name] = true
	}
	for _, a := range jobArguments.Artifacts {
		argArtifactNames[a.Name] = true
	}
	// Pass arguments to the DAG task if referenced by its inputs/outputs
	ioParameters := append(inputs.Parameters, outputs.Parameters...)
	for _, p := range ioParameters {
		if _, found := argParameterNames[p.Name]; found {
			value := fmt.Sprintf("{{inputs.parameters.%s}}", p.Name)
			dagTask.Arguments.Parameters = append(
				dagTask.Arguments.Parameters, argov1alpha1.Parameter{
					Name:  p.Name,
					Value: &value,
				})
		}
	}
	ioArtifacts := append(inputs.Artifacts, outputs.Artifacts...)
	for _, a := range ioArtifacts {
		if _, found := argArtifactNames[a.Name]; found {
			dagTask.Arguments.Artifacts = append(
				dagTask.Arguments.Artifacts, argov1alpha1.Artifact{
					Name: a.Name,
					From: fmt.Sprintf("{{inputs.artifacts.%s}}", a.Name),
				})
		}
	}

	return dagTask
}

func buildDAGTemplate(
	name string,
	dagTasks []argov1alpha1.DAGTask,
	jobArguments argov1alpha1.Arguments) argov1alpha1.Template {

	// Pass all arguments as inputs to the DAG template.
	inputs := argov1alpha1.Inputs{}
	for _, p := range jobArguments.Parameters {
		inputs.Parameters = append(inputs.Parameters, argov1alpha1.Parameter{Name: p.Name})
	}
	for _, a := range jobArguments.Artifacts {
		inputs.Artifacts = append(inputs.Artifacts, argov1alpha1.Artifact{Name: a.Name})
	}

	template := argov1alpha1.Template{
		Name: "kubebench-job-workflow-entrypoint",
		DAG: &argov1alpha1.DAGTemplate{
			Tasks: dagTasks,
		},
		Inputs: inputs,
	}

	return template
}

func buildContainerTemplate(
	templateName string,
	container *corev1.Container,
	wfInfo *workflowInfo,
	inputs argov1alpha1.Inputs,
	outputs argov1alpha1.Outputs) argov1alpha1.Template {

	modSpec := &mod.ResourceModSpec{
		VolumeMounts: wfInfo.managedVolumeMounts,
		Env:          wfInfo.env,
	}
	modContainer := mod.ModifyContainerV1(*container, modSpec)

	template := argov1alpha1.Template{
		Name:      templateName,
		Container: &modContainer,
	}

	template.Inputs = inputs
	template.Outputs = outputs

	return template
}

// The "IO" template is useful for inputs/outputs that do not require additional operations.
func buildIOTemplate(
	name string,
	wfaContainer *corev1.Container,
	wfInfo *workflowInfo,
	inputs argov1alpha1.Inputs,
	outputs argov1alpha1.Outputs) argov1alpha1.Template {

	wfaContainer.Command = []string{"sleep", "1"}
	template := buildContainerTemplate(name, wfaContainer, wfInfo, inputs, outputs)
	return template
}

func buildResourceConfigTemplate(
	name string,
	wfaContainer *corev1.Container,
	resSpec *kbjobv1alpha2.ResourceSpec,
	wfInfo *workflowInfo,
	inputs argov1alpha1.Inputs) argov1alpha1.Template {

	templateName := getName(configuratorTemplateNameFmt, name)
	outputName := getName(configuratorOutputNameFmt, name)

	confInputStr := buildConfiguratorInputStr(resSpec, wfInfo)
	outputFile := path.Join(
		fmt.Sprintf(constants.WorkflowExpPathFmt, wfInfo.experimentID), outputName)
	wfaContainer.Command = []string{
		"configurator",
		"--input-params", confInputStr,
		"--output-file", outputFile,
	}

	outputs := argov1alpha1.Outputs{
		Parameters: []argov1alpha1.Parameter{
			{
				Name:       outputName,
				GlobalName: outputName,
				ValueFrom: &argov1alpha1.ValueFrom{
					Path: outputFile,
				},
			},
		},
	}

	template := buildContainerTemplate(templateName, wfaContainer, wfInfo, inputs, outputs)

	return template
}

func buildResourceCreateTemplate(
	name string,
	wfaContainer *corev1.Container,
	resSpec *kbjobv1alpha2.ResourceSpec,
	wfInfo *workflowInfo) argov1alpha1.Template {

	templateName := getName(resourceCreateTemplateNameFmt, name)
	inputName := getName(configuratorOutputNameFmt, name)
	outputName := getName(resourceCreateOutputNameFmt, name)

	outputFile := path.Join(
		fmt.Sprintf(constants.WorkflowExpPathFmt, wfInfo.experimentID), outputName)

	var numCopies int
	if resSpec.Options != nil && resSpec.Options.NumCopies > 0 {
		numCopies = resSpec.Options.NumCopies
	} else {
		numCopies = 1
	}
	wfaContainer.Command = []string{
		"resource-manager",
		"--action", "create",
		"--num-copies", strconv.Itoa(numCopies),
		"--input-data", fmt.Sprintf("{{workflow.outputs.parameters.%s}}", inputName),
		"--output-file", outputFile,
	}

	inputs := argov1alpha1.Inputs{}
	outputs := argov1alpha1.Outputs{
		Parameters: []argov1alpha1.Parameter{
			{
				Name:       outputName,
				GlobalName: outputName,
				ValueFrom: &argov1alpha1.ValueFrom{
					Path: outputFile,
				},
			},
		},
	}

	template := buildContainerTemplate(templateName, wfaContainer, wfInfo, inputs, outputs)

	return template
}

func buildResourceAutoWatchTemplate(
	name string,
	wfaContainer *corev1.Container,
	resSpec *kbjobv1alpha2.ResourceSpec,
	wfInfo *workflowInfo) argov1alpha1.Template {

	templateName := getName(resourceAutoWatchTemplateNameFmt, name)
	inputName := getName(resourceCreateOutputNameFmt, name)
	outputName := getName(resourceAutoWatchOutputNameFmt, name)

	outputFile := path.Join(
		fmt.Sprintf(constants.WorkflowExpPathFmt, wfInfo.experimentID), outputName)

	var timeout string
	if resSpec.Options != nil && resSpec.Options.AutoWatch != nil {
		timeout = resSpec.Options.AutoWatch.Timeout
	}
	wfaContainer.Command = []string{
		"resource-manager",
		"--action", "auto-watch",
		"--timeout", timeout,
		"--input-data", fmt.Sprintf("{{workflow.outputs.parameters.%s}}", inputName),
		"--output-file", outputFile,
	}

	inputs := argov1alpha1.Inputs{}
	outputs := argov1alpha1.Outputs{
		Parameters: []argov1alpha1.Parameter{
			{
				Name:       outputName,
				GlobalName: outputName,
				ValueFrom: &argov1alpha1.ValueFrom{
					Path: outputFile,
				},
			},
		},
	}

	template := buildContainerTemplate(templateName, wfaContainer, wfInfo, inputs, outputs)

	return template
}

func buildConfiguratorInputStr(
	resSpec *kbjobv1alpha2.ResourceSpec,
	wfInfo *workflowInfo) string {

	// Generate manifest generation spec
	manifestGenSpec := wfacommon.ManifestGenSpec{
		Manifest:     resSpec.Manifest,
		ManifestFrom: resSpec.ManifestFrom,
	}

	// Generate volume info
	volsToMnt := []corev1.Volume{}
	managedVolsToMnt := []corev1.Volume{}
	for i, vm := range resSpec.VolumeMounts {
		// Add volumes to be mounted
		if v, found := wfInfo.volumeMap[vm.Name]; found {
			volsToMnt = append(volsToMnt, v)
		}
		// Add managed volumes to be mounted
		// Detect if managed volume is explicitly mounted, and change subpath if so.
		if v, found := wfInfo.managedVolumeMap[vm.Name]; found {
			managedVolsToMnt = append(managedVolsToMnt, v)
			subPath := resSpec.VolumeMounts[i].SubPath
			if subPath == "" {
				subPath = wfInfo.experimentID
			} else {
				subPath = wfInfo.experimentID + "/" + subPath
			}
			resSpec.VolumeMounts[i].SubPath = subPath
		}
	}
	allVolsToMnt := []corev1.Volume{}
	if resSpec.Options != nil && resSpec.Options.MountManagedVolumes {
		allVolsToMnt = append(volsToMnt, wfInfo.managedVolumes...)
	} else {
		allVolsToMnt = append(volsToMnt, managedVolsToMnt...)
	}

	// Generate manifest modification spec
	manifestModSpec := wfacommon.ManifestModSpec(
		mod.ResourceModSpec{
			Namespace:       wfInfo.namespace,
			OwnerReferences: wfInfo.ownerReferences,
			Labels:          wfInfo.labels,
			Volumes:         allVolsToMnt,
			VolumeMounts:    append(resSpec.VolumeMounts, wfInfo.managedVolumeMounts...),
			Env:             wfInfo.env,
		},
	)

	// Generate configurator input in string form
	confInput := &wfacommon.ConfiguratorInput{
		ManifestGenSpec: &manifestGenSpec,
		ManifestModSpec: &manifestModSpec,
	}
	confInputByte, _ := json.Marshal(confInput)
	confInputStr := string(confInputByte)
	return confInputStr
}
