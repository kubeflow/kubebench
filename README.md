# kubebench

The goal of Kubebench is to make it easy to run benchmark jobs on [Kubeflow](https://github.com/kubeflow/kubeflow) with various system and model settings. Kubebench enables benchmarks by leveraging Kubeflow's capability of managing TFJobs, as well as [Argo](https://github.com/argoproj/argo) based workflows.


## Getting Started

### Prerequisites

  - Kubernetes >= 1.8
  - Ksonnet >= 0.10

### Installation

  - Install Dependencies (Kubebench depends on an existing Kubeflow deployment. For details about using Kubeflow, please refer to [Kubeflow user guide](https://github.com/kubeflow/kubeflow/blob/master/user_guide.md))

    ```bash
    KF_VERSION=master
    KF_ENV=local
    NAMESPACE=default

    # Initialize Ksonnet app
    ks init my-kubeflow
    cd my-kubeflow

    # Install required Kubeflow packages
    ks registry add kubeflow github.com/kubeflow/kubeflow/tree/${KF_VERSION}/kubeflow
    ks pkg install kubeflow/core@${KF_VERSION}
    ks pkg install kubeflow/tf-job@${KF_VERSION}
    ks pkg install kubeflow/argo@${KF_VERSION}

    # Create components
    ks generate core kubeflow-core --name=kubeflow-core
    ks generate argo kubeflow-argo --name=kubeflow-argo

    # Configure environment
    ks env add ${KF_ENV}
    ks env set ${KF_ENV} --namespace ${NAMESPACE}

    # Deploy the components
    ks apply ${KF_ENV} -c kubeflow-core
    ks apply ${KF_ENV} -c kubeflow-argo

    # Configure service account to grant Argo more privileges
    kubectl create rolebinding default-admin --clusterrole=admin --serviceaccount=default:default
    ```

  - Install Kubebench

    ```bash
    KB_VERSION=master
    ks registry add kubebench github.com/kubeflow/kubebench/tree/${KB_VERSION}/kubebench
    ks pkg install kubebench/kubebench-job@${KB_VERSION}
    ```

  - Create a persistent volume claim for data storage
    - currently this is the only supported way to store benchmark configurations and results. In the future this will be simplified, and more options of configuration/result storage will be provided.
    - the persistent volume needs to have ReadWriteMany access mode.
    - we provide an example PVC setup below based on NFS, the example config file is in [examples/tf_cnn_benchmarks directory](https://github.com/kubeflow/kubebench/blob/master/examples/tf_cnn_benchmarks/nfs_pvc.yaml)
    - assuming that you have NFS setup, edit the example config file with the right server address and exported path of your NFS server, then run the following command

    ```bash
    kubectl create -f nfs_pvc.yaml --namespace ${NAMESPACE}
    ```

### Run a Kubebench Job

  - Copy benchmark configuration file to the persistent volume
    - an example benchmark config is in [examples/tf_cnn_benchmarks directory](https://github.com/kubeflow/kubebench/blob/master/examples/tf_cnn_benchmarks/job_config.yaml)
    - assuming that your NFS exported path is `/var/nfs/kubebench`, as is in the PVC config example, run the following command on the NFS server

    ```
    mkdir -p /var/nfs/kubebench/config
    cp job_config.yaml /var/nfs/kubebench/config
    ```

  - Generate, configure, and deploy a kubebench-job

    ```
    CONFIG_NAME="job_config"
    JOB_NAME="my-benchmark"
    PVC_NAME="kubebench-pvc"
    PVC_MOUNT="/kubebench"

    ks generate kubebench-job ${JOB_NAME} --name=${JOB_NAME}

    ks param set ${JOB_NAME} name ${CONFIG_NAME}
    ks param set ${JOB_NAME} namespace ${NAMESPACE}
    ks param set ${JOB_NAME} config_image gcr.io/xyhuang-kubeflow/kubebench-configurator:v20180522-1
    ks param set ${JOB_NAME} report_image gcr.io/xyhuang-kubeflow/kubebench-tf-cnn-csv-reporter:v20180522-1
    ks param set ${JOB_NAME} config_args -- --config-file=${PVC_MOUNT}/config/${CONFIG_NAME}.yaml
    ks param set ${JOB_NAME} report_args -- --output-file=${PVC_MOUNT}/output/results.csv
    ks param set ${JOB_NAME} pvc_name ${PVC_NAME}
    ks param set ${JOB_NAME} pvc_mount ${PVC_MOUNT}

    ks apply ${KF_ENV} -c ${JOB_NAME}
    ```

  - Track the status of your job

    ```
    kubectl get -o yaml workflows ${JOB_NAME}
    ```

  - Once the job is finished, you can find the results under your specified output directory in the NFS, if you used the same configurations as is in this example, the output directory will be in `/var/nfs/kubebench/output`.

## Design Document

For additional information about motivation and design for this project please refer to [kubebench_design.md](./doc/kubebench_design.md)
