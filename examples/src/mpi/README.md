# Launch MPI Job benchmark using kubebench
This tutorial will show an example to launching OpenMPI benchmark job on [Kubeflow](https://github.com/kubeflow/kubeflow) 

## Quick Start

### Prerequisites

  - Kubernetes >= 1.10 
    - GPU nodes with [Nvidia Device Plugin](https://github.com/NVIDIA/k8s-device-plugin)
  - Ksonnet >= 0.11
  - Kubeflow master
    - Required modules: argo, mpi-operator
  - For the quick-starter installation, Kubernetes nodes need to support NFS mounting


### Installation

- Init Project

  ```
  ks init ${BENCHMARK_PROJECT} && cd ${BENCHMARK_PROJECT}
  ks registry add kubeflow github.com/kubeflow/kubeflow/tree/master/kubeflow
  ```

- Setup NFS used by Kubebench

  Please check [page](https://github.com/kubeflow/kubebench/blob/master/README.md#installation) to install Kubebench quick-starter package

- Setup required components
  ```
  # Install dependency packages 
  ks pkg install kubeflow/common@master
  ks pkg install kubeflow/argo@master
  ks pkg install kubeflow/mpi-job@master
  ks pkg install kubeflow/kubebench@master
  
  # Generate Manifests
  ks generate argo argo
  ks generate mpi-operator mpi-operator
  
  # Customize your deployment
  ks param set mpi-operator image seedjeffwan/mpi-operator:latest
  
  # Deploy required components
  Ks apply default
  ```
  > Note: Default mpi-operator image doesn't have latest [change](https://github.com/kubeflow/mpi-operator/pull/89). That's why we use self built image instead.

### Run a Kubebench Job

- Configure Kubebench job
  ```
  ks generate kubebench-job ${JOB_NAME}
     
  ks param set ${JOB_NAME} mainJobConfig mpi/mpi-job-dummy.yaml
  ks param set ${JOB_NAME} mainJobKsPackage mpi-job
  ks param set ${JOB_NAME} mainJobKsPrototype mpi-job-custom
  ks param set ${JOB_NAME} mainJobKsRegistry github.com/kubeflow/kubeflow/tree/master/kubeflow
     
  ks param set ${JOB_NAME} controllerImage seedjeffwan/configurator:latest
  ks param set ${JOB_NAME} postJobImage seedjeffwan/mpi-post-processor:latest
  
  # Optional
  ks param set ${JOB_NAME} githubTokenSecret github-token
  ks param set ${JOB_NAME} githubTokenSecretKey GITHUB_TOKEN
  ```

  Note:
  * MPI Job configuration file [mpi/mpi-job-dummy.yaml](../../config/mpi/mpi-job-dummy.yaml) is already in your NFS config folder.  
  * Without github token, you might run into API rate limits. Generate one at [Github Token](https://github.com/settings/tokens) and create a kubernetes secret file.
   
  ```
  apiVersion: v1
  kind: Secret
  metadata:
    name: github-token
  type: Opaque
  data:
    GITHUB_TOKEN: YOUR_BASE64_ENCODED_TOKEN
  ```

- Launch MPI Training job
  ```
  ks apply default -c ${JOB_NAME}
  ```

## FAQ

#### How can I look at benchmark workflow in Argo UI? 

Argo UI by default is backed by Amabassdor that we didn't install in this tutorial. We need to make minor change to expose Argo UI to frontend. 

Edit Argo UI deployment resource file by `kubectl edit deployment argo-ui`. Look for `BASE_HREF`, change environment value from `/argo/` to `/`.

Wait for the new pod start, then run `kubectl port-forward deployment/argo-ui 8001:8001` to proxy port. Now you can visit `localhost:8001` for workflow details. 


#### What does the output looks like? 

`kubectl port-forward deployment/kubebench-nfs-deploy 8000:80` and open `localhost:8000` to check experiment outputs. 

```
|-- config
|   `-- mpi
|       `-- mpi-job-dummy.yaml
|
|-- data
`-- experiments
    |-- mpi-job-dummy-201902170707-8i0d
    |   |-- config
    |   |   |-- kf-job-manifest.yaml
    |   |   `-- mpi-job-dummy.yaml
    |   |-- output
    |   |   `-- mpi-job-dummy-201902170707-8i0d-launcher-82mlc
    |   `-- result
    |       `-- result.json
    `-- report.csv
```
