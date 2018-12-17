# kubebench operator

The goal of kubebench-operator is to facilitate running kubebench workflows. It is the high-level wrapper of Argo and KubebenchJob. To submit a KubebenchJob a workflow yaml file is needed.

## Prerequisites

  - Kubernetes >= 1.8
  - [Ksonnet](ksonnet.io) >= 0.11
  - [Kubeflow](https://github.com/kubeflow/kubeflow/) >= 0.3
    - Required modules: [Argo](https://github.com/argoproj/argo), [tf-operator](https://github.com/kubeflow/tf-operator)

## Description

Kubebench operator controlls object of type ```KubebenchJob```. A workflow example can be below:

```
apiVersion: kubebench.operator/v1
kind: KubebenchJob
metadata:
  name: kubebench-job
  namespace: default
spec:
  serviceAccount: default
  volumeSpecs:
    configVolume:
      name: my-config-volume
      persistentVolumeClaim:
        claimName: kubebench-config-pvc
    experimentVolume:
      name: my-experiment-volume
      persistentVolumeClaim:
        claimName: kubebench-exp-pvc
  secretSpecs: # optional
    githubTokenSecret: # optional
      secretName: my-github-token-secret
      secretKey: my-github-token-secret-key
    gcpCredentialsSecret: # optional
      secretName: my-gcp-credentials-secret
      secretKey: my-gcp-credentials-secret-key
  jobSpecs:
    preJob: # optional
      container: # optional between "container" and "resource"
        name: my-prejob
        image: gcr.io/myprejob-image:latest # change it before using
    mainJob: # mandatory
      resource: # optional between "container" and "resource"
        manifestTemplate:
          valueFrom:
            ksonnet: # optional, more types in the future
              prototype: kubebench-example-tfcnn
              package: kubebench-examples
              registry: github.com/kubeflow/kubebench/tree/master/kubebench
        manifestParameters:
          valueFrom:
            path: abc/def/ghi.yaml
        createSuccessCondition: createSuccess # optional
        createFailureCondition: createFailure # optional
        runSuccessCondition: runSuccess # optional
        runFailureCondition: runFailre # optional
        #other optional fields: "manifest" - string of raw manifest
    postJob: # optional
      container: # optional between "container" and "resource"
        name: my-postjob
        image: gcr.io/kubeflow-images-public/kubebench/kubebench-example-tf-cnn-post-processor:3c75b50
  reportSpecs: # optional
    csv: # optional
        - inputPath: result.json
          outputPath: report.csv
```

## Quickstart

In order to run kubebench-operator be sure to have all the prerequisites deployed/installed. 
Deploy CRD itself first:

```
kubectl apply -f examples/crd.yaml
```

Then you have to deploy controller itself:

```
kubectl apply -f examples/crd-deployment.yaml
```

And finally to allow your serviceAccount to create pods, modify and submit ```cluster-role.yaml``` :

```
kubectl apply -f examples/cluster-role.yaml
```

## Development

For the information on developing visit: [developer guide](https://github.com/kubeflow/kubebench/controller/kubebench-operator/DEVELOPER.md)