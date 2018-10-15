package util

import (
	"fmt"
	"strings"

	workflow_v1 "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	kubebenchjob "github.com/kubeflow/kubebench/controller/kubebench-operator/pkg/apis/kubebenchjob/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	kubebenchConfigVol      = "kubebench-config-volume"
	kubebenchDataVol        = "kubebench-data-volume"
	kubebenchExpVol         = "kubebench-exp-volume"
	kubebenchGithubTokenVol = "kubebench-github-token-volume"
	kubebenchGcpCredsVol    = "kubebench-gcp-credentials-volume"
	kubebenchConfigRoot     = "/kubebench/config"
	kubebenchDataRoot       = "/kubebench/data"
	kubebenchExpRoot        = "/kubebench/experiments"
	configuratorOutputDir   = "/kubebench/configurator/output"
	manifestOutput          = configuratorOutputDir + "/kf-job-manifest.yaml"
	experimentIdOutput      = configuratorOutputDir + "/experiment-id"
)

func ConvertVolumesToString(volumes []apiv1.Volume) string {
	result := []string{}
	for _, volume := range volumes {
		result = append(result, fmt.Sprintf("{\"name\": \"%s\", \"persistentVolumeClaim\": {\"claimName\": \"%s\"}}", volume.Name, volume.VolumeSource.PersistentVolumeClaim.ClaimName))
	}
	tempResult := strings.Join(result, ",")
	return "[" + tempResult + "]"
}

func ConvertVolumeMountsToString(vMounts []apiv1.VolumeMount) string {
	result := []string{}
	for _, volMnt := range vMounts {
		result = append(result, fmt.Sprintf("{\"mountPath\": \"%s\", \"name\": \"%s\"}", volMnt.MountPath, volMnt.Name))
	}

	tempResult := strings.Join(result, ",")
	return "[" + tempResult + "]"
}

func ConvertEnvVarsToString(envs []apiv1.EnvVar) string {
	result := []string{}
	for _, env := range envs {
		result = append(result, fmt.Sprintf("{\"name\": \"%s\", \"value\": \"%s\"}", env.Name, env.Value))
	}

	tempResult := strings.Join(result, ",")
	return "[" + tempResult + "]"
}

// think over
func CreateKeyValuePairs(m map[string]string) string {
	result := []string{}
	for key, value := range m {
		result = append(result, fmt.Sprintf("\"%s\": \"%s\"", key, value))
	}
	return "{" + strings.Join(result, ",") + "}"
}

func BuildTemplate(stepName string, image string, command []string, envVars []apiv1.EnvVar, volMnts []apiv1.VolumeMount, inParams workflow_v1.Inputs, outParams workflow_v1.Outputs) workflow_v1.Template {
	template := workflow_v1.Template{
		Name: stepName,
		Container: &apiv1.Container{
			Image:        image,
			VolumeMounts: volMnts,
			Env:          envVars,
		},
	}
	if len(command) > 0 {
		template.Container.Command = command
	}

	if len(inParams.Parameters) > 0 {
		template.Inputs = inParams
	}

	if len(outParams.Parameters) > 0 {
		template.Outputs = outParams
	}

	return template
}

func BuildResourceTemplate(stepName string, action string, manifest string, successCondition string, failureCondition string, inParams workflow_v1.Inputs, outParams workflow_v1.Outputs) workflow_v1.Template {
	template := workflow_v1.Template{
		Name: stepName,
		Resource: &workflow_v1.ResourceTemplate{
			Action:   action,
			Manifest: manifest,
		},
	}

	if successCondition != "" {
		template.Resource.SuccessCondition = successCondition
	}

	if failureCondition != "" {
		template.Resource.FailureCondition = failureCondition
	}

	if len(inParams.Parameters) > 0 {
		template.Inputs = inParams
	}

	if len(outParams.Parameters) > 0 {
		template.Outputs = outParams
	}
	return template
}

func BuildStep(stepName string, template string, argParams workflow_v1.Arguments) workflow_v1.WorkflowStep {
	step := workflow_v1.WorkflowStep{
		Name:     stepName,
		Template: template,
	}

	if len(argParams.Parameters) > 0 {
		step.Arguments = argParams
	}

	return step
}

func expEnvVars(isConfigurator bool) []apiv1.EnvVar {
	parameters := make([]apiv1.EnvVar, 0)
	if isConfigurator {
		parameters = append(parameters, apiv1.EnvVar{
			Name:  "KUBEBENCH_EXP_ID",
			Value: "null",
		})
	} else {
		parameters = append(parameters,
			apiv1.EnvVar{
				Name:  "KUBEBENCH_EXP_ID",
				Value: "{{inputs.parameters.experiment-id}}",
			},
		)
	}

	parameters = append(parameters,
		apiv1.EnvVar{
			Name:  "KUBEBENCH_EXP_PATH",
			Value: "$(KUBEBENCH_EXP_ROOT)/$(KUBEBENCH_EXP_ID)",
		},
		apiv1.EnvVar{
			Name:  "KUBEBENCH_EXP_CONFIG_PATH",
			Value: "$(KUBEBENCH_EXP_PATH)/config",
		},
		apiv1.EnvVar{
			Name:  "KUBEBENCH_EXP_OUTPUT_PATH",
			Value: "$(KUBEBENCH_EXP_PATH)/output",
		},
		apiv1.EnvVar{
			Name:  "KUBEBENCH_EXP_RESULT_PATH",
			Value: "$(KUBEBENCH_EXP_PATH)/result",
		})
	return parameters
}

func CreateArgumentsForString(names []string, values []string) workflow_v1.Arguments {
	newArguments := workflow_v1.Arguments{}

	params := []workflow_v1.Parameter{}
	for i := 0; i < len(names); i++ {
		params = append(params, workflow_v1.Parameter{
			Name:  names[i],
			Value: &values[i],
		})
	}

	newArguments.Parameters = params
	return newArguments
}

func CreateInputs(names []string) workflow_v1.Inputs {
	inputs := workflow_v1.Inputs{}
	params := []workflow_v1.Parameter{}
	for i := 0; i < len(names); i++ {
		params = append(params, workflow_v1.Parameter{
			Name: names[i],
		})
	}

	if len(params) > 0 {
		inputs.Parameters = params
	}
	return inputs
}

func CreateOutputs(names []string, paths []string) workflow_v1.Outputs {
	outputs := workflow_v1.Outputs{}
	params := []workflow_v1.Parameter{}
	for i := 0; i < len(names); i++ {
		params = append(params, workflow_v1.Parameter{
			Name: names[i],
			ValueFrom: &workflow_v1.ValueFrom{
				Path: paths[i],
			},
		})
	}
	if len(params) > 0 {
		outputs.Parameters = params
	}

	return outputs
}

func ConvertKubebenchJobToArgoWorkflow(kbJob *kubebenchjob.KubebenchJob) (wkflw *workflow_v1.Workflow, err error) {

	//get rid of hardcoded fields, update types + yaml file

	reporterArgs := make([]string, 0)
	jbSpec := kbJob.Spec

	postJobArgs := []string{}
	reporterType := "csv"

	mainJobKsPrototypeRef := jbSpec.Jobs.MainJob.Resource.ManifestTemplate.ValueFrom.KsonnetSpec
	mainJobConfig := "tf-cnn/tf-cnn-dummy.yaml"

	ownerReferences := map[string]string{
		"apiVersion": "argoproj.io/v1alpha1",
		//		"blockOwnerDeletion": "true",
		"kind": "Workflow",
		"name": "{{workflow.name}}",
		"uid":  "{{workflow.uid}}",
	}

	if len(jbSpec.Report.CSV) != 0 {
		reporterArgs = append(reporterArgs, "--input-file="+jbSpec.Report.CSV[0]["inputPath"])
		reporterArgs = append(reporterArgs, "--output-file="+jbSpec.Report.CSV[0]["outputPath"])
	}

	secretEnvVars := make([]apiv1.EnvVar, 0)
	secretVols := make([]apiv1.Volume, 0)
	secretVolMnts := make([]apiv1.VolumeMount, 0)

	if len(jbSpec.Secrets.GithubToken) != 0 {

		secretEnvVars = append(secretEnvVars, apiv1.EnvVar{
			Name: "GITHUB_TOKEN",
			ValueFrom: &apiv1.EnvVarSource{
				SecretKeyRef: &apiv1.SecretKeySelector{
					LocalObjectReference: apiv1.LocalObjectReference{
						Name: jbSpec.Secrets.GithubToken["secretName"],
					}, //fix!!
					Key: jbSpec.Secrets.GithubToken["secretKey"],
				},
			},
		})

		secretVols = append(secretVols, apiv1.Volume{
			Name: kubebenchGithubTokenVol,
			VolumeSource: apiv1.VolumeSource{
				Secret: &apiv1.SecretVolumeSource{
					SecretName: jbSpec.Secrets.GithubToken["secretName"],
				},
			},
		})

		secretVolMnts = append(secretVolMnts, apiv1.VolumeMount{
			Name:      kubebenchGithubTokenVol,
			MountPath: "/secret/github-token",
		})
	}

	if len(jbSpec.Secrets.GCPCredentials) != 0 {

		secretEnvVars = append(secretEnvVars, apiv1.EnvVar{
			Name:  "GOOGLE_APPLICATION_CREDENTIALS",
			Value: "/secret/gcp-credentials/" + jbSpec.Secrets.GCPCredentials["secretKey"],
		})

		secretVols = append(secretVols, apiv1.Volume{
			Name: kubebenchGcpCredsVol,
			VolumeSource: apiv1.VolumeSource{
				Secret: &apiv1.SecretVolumeSource{
					SecretName: jbSpec.Secrets.GCPCredentials["secretName"],
				},
			},
		})

		secretVolMnts = append(secretVolMnts, apiv1.VolumeMount{
			Name:      kubebenchGcpCredsVol,
			MountPath: "/secret/gcp-credentials",
		})
	}

	baseEnvVars := []apiv1.EnvVar{
		apiv1.EnvVar{
			Name:  "KUBEBENCH_CONFIG_ROOT",
			Value: kubebenchConfigRoot,
		},
		apiv1.EnvVar{
			Name:  "KUBEBENCH_EXP_ROOT",
			Value: kubebenchExpRoot,
		},
		apiv1.EnvVar{
			Name:  "KUBEBENCH_DATA_ROOT",
			Value: kubebenchDataRoot,
		},
	}

	baseVols := []apiv1.Volume{
		apiv1.Volume{
			Name: kubebenchConfigVol,
			VolumeSource: apiv1.VolumeSource{
				PersistentVolumeClaim: &apiv1.PersistentVolumeClaimVolumeSource{
					ClaimName: jbSpec.Volumes.ConfigVolume.PersistentVolumeClaim.ClaimName,
				},
			},
		},
		apiv1.Volume{
			Name: kubebenchExpVol,
			VolumeSource: apiv1.VolumeSource{
				PersistentVolumeClaim: &apiv1.PersistentVolumeClaimVolumeSource{
					ClaimName: jbSpec.Volumes.ExperimentVolume.PersistentVolumeClaim.ClaimName,
				},
			},
		},
	}

	baseVolMnts := []apiv1.VolumeMount{
		apiv1.VolumeMount{
			Name:      kubebenchConfigVol,
			MountPath: kubebenchConfigRoot,
		},
		apiv1.VolumeMount{
			Name:      kubebenchExpVol,
			MountPath: kubebenchExpRoot,
		},
	}

	configuratorCommand := []string{
		"configurator",
		"--template-ref=" + CreateKeyValuePairs(mainJobKsPrototypeRef),
		"--config=" + mainJobConfig,
		"--namespace=" + kbJob.ObjectMeta.Namespace,
		"--owner-references=" + "[" + CreateKeyValuePairs(ownerReferences) + "]",
		"--volumes=" + ConvertVolumesToString(baseVols),
		"--volume-mounts=" + ConvertVolumeMountsToString(baseVolMnts),
		"--env-vars=" + ConvertEnvVarsToString(append(baseEnvVars, expEnvVars(true)...)),
		"--manifest-output=" + manifestOutput,
		"--experiment-id-output=" + experimentIdOutput,
	}

	//controllerImage := "gcr.io/xyhuang-kubeflow/kubebench-controller:v0.3.0"
	controllerImage := "gcr.io/kubeflow-images-public/kubebench/kubebench-controller:3c75b50"

	result := &workflow_v1.Workflow{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "argoproj.io/v1alpha1",
			Kind:       "Workflow",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      kbJob.ObjectMeta.Name,
			Namespace: kbJob.ObjectMeta.Namespace,
		},
		Spec: workflow_v1.WorkflowSpec{
			ServiceAccountName: jbSpec.ServiceAccount,
			Entrypoint:         "kubebench-workflow",
			Volumes:            append(baseVols, secretVols...),
			Templates: []workflow_v1.Template{
				workflow_v1.Template{
					Name: "kubebench-workflow",
					Steps: [][]workflow_v1.WorkflowStep{
						[]workflow_v1.WorkflowStep{
							BuildStep("run-configurator", "configurator", CreateArgumentsForString([]string{}, []string{})),
						},
						[]workflow_v1.WorkflowStep{
							BuildStep("launch-main-job", "main-job", CreateArgumentsForString(
								[]string{"kf-job-manifest", "experiment-id"}, []string{"{{steps.run-configurator.outputs.parameters.kf-job-manifest}}", "{{steps.run-configurator.outputs.parameters.experiment-id}}"})),
						},
						[]workflow_v1.WorkflowStep{
							BuildStep("wait-for-main-job",
								"main-job-monitor", CreateArgumentsForString([]string{"kf-job-manifest"}, []string{"{{steps.run-configurator.outputs.parameters.kf-job-manifest}}"})),
						},
						[]workflow_v1.WorkflowStep{
							BuildStep("run-post-job", "post-job", CreateArgumentsForString([]string{"kf-job-manifest", "experiment-id"}, []string{"{{steps.run-configurator.outputs.parameters.kf-job-manifest}}", "{{steps.run-configurator.outputs.parameters.experiment-id}}"})),
						},
						[]workflow_v1.WorkflowStep{
							BuildStep("run-reporter", "reporter", CreateArgumentsForString([]string{"kf-job-manifest", "experiment-id"}, []string{"{{steps.run-configurator.outputs.parameters.kf-job-manifest}}", "{{steps.run-configurator.outputs.parameters.experiment-id}}"})),
						},
					},
				},

				BuildTemplate("configurator", controllerImage, configuratorCommand, append(secretEnvVars, baseEnvVars...), append(secretVolMnts, baseVolMnts...), CreateInputs([]string{}), CreateOutputs([]string{"kf-job-manifest", "experiment-id"}, []string{manifestOutput, experimentIdOutput})),

				BuildResourceTemplate("main-job", "create", "{{inputs.parameters.kf-job-manifest}}", "status.startTime", "", CreateInputs([]string{"kf-job-manifest"}), CreateOutputs([]string{}, []string{})),

				BuildResourceTemplate("main-job-monitor", "get", "{{inputs.parameters.kf-job-manifest}}", "status.completionTime", "", CreateInputs([]string{"kf-job-manifest"}), CreateOutputs([]string{}, []string{})),

				BuildTemplate("post-job", jbSpec.Jobs.PostJob.Container.Image, postJobArgs, append(baseEnvVars, expEnvVars(false)...), baseVolMnts, CreateInputs([]string{"experiment-id"}), CreateOutputs([]string{}, []string{})),

				BuildTemplate("reporter", controllerImage, append([]string{"reporter", reporterType}, reporterArgs...), append(append(secretEnvVars, baseEnvVars...), expEnvVars(false)...), append(secretVolMnts, baseVolMnts...), CreateInputs([]string{"experiment-id"}), CreateOutputs([]string{}, []string{})),
			},
		},
	}

	//controller image inside yaml ?

	return result, nil
}

func GetJobName(kbJob *kubebenchjob.KubebenchJob) (name string) {
	return kbJob.ObjectMeta.Name
}
