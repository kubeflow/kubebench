# kubebench

The goal of Kubebench is to make it easy to run benchmark jobs on [Kubeflow](https://github.com/kubeflow/kubeflow) with various system and model settings. Kubebench enables benchmarks by leveraging Kubeflow's capability of managing TFJobs, as well as [Argo](https://github.com/argoproj/argo) based workflows.


## Quick Start

NOTE: the quick start guide serves as a demo that helps you quickly go through a Kubebench Job. The components installed may not be suitable for production use. Please refer to detailed user guide for proper configuration of Kubebench Jobs.

### Prerequisites

  - Kubernetes >= 1.8
  - Ksonnet >= 0.11
  - Kubeflow >= 0.3
    - Required modules: argo, tf-operator
  - For the quick-starter installation, Kubernetes nodes need to support NFS mounting

### Installation

  - Install Dependencies (Kubebench depends on an existing Kubeflow deployment. For details about using Kubeflow, please refer to [Kubeflow documentation](https://www.kubeflow.org/docs/started/getting-started/))

  - Install Kubebench quick-starter package

    In your Ksonnet app root, run the followings:

    ```bash
    export KB_VERSION=master
    export KB_ENV=default
    source ../env.sh
    
    curl https://raw.githubusercontent.com/kubeflow/kubebench/master/scripts/install_quickstarter.sh | bash
    ```

  - View the Kubebench directory contents

    The installer comes with a simple file server that allows you to view the contents of Kubebench directory through browser. You may find details of the file server service through:

    ```bash
    kubectl get svc kubebench-nfs-file-server-svc -o wide
    ```

    Alternatively, you can also access the deployed NFS service directly. You may find details of the nfs service through:

    ```bash
    kubectl get svc kubebench-nfs-svc -o wide
    ```

### Run a Kubebench Job

  - Create a kubebench-job

    ```bash
    JOB_NAME="my-benchmark"

    ks pkg install kubebench/kubebench-job@${KB_VERSION}
    ks generate kubebench-job ${JOB_NAME}

    ks apply ${KB_ENV} -c ${JOB_NAME}
    ```

  - Track the status of your job

    The Kubebench Job will be deployed as an Argo Workflow, you may go to Argo dashboard to track the progress of the job.

    Alternatively, you can also use the followings in command line:

    ```bash
    kubectl get -o yaml workflows ${JOB_NAME}
    ```

### View results

  - Once the job is finished, you can find the results under the experiment directory in the NFS, the details of the particular experiment is located at `/experiments/<EXPERIMENT_UID>`. You may also see a csv file at `/experiments/report.csv`, if you run multiple experiments, the aggregated results will be recorded here.

### Cleanups

  - Delete the kubebench-job

    ```bash
    ks delete ${KB_ENV} -c ${JOB_NAME}
    ```

  - Uninstall quickstarter

    ```bash
    curl https://raw.githubusercontent.com/kubeflow/kubebench/master/scripts/uninstall_quickstarter.sh | bash
    ```

## Design Document

For additional information about motivation and design for this project please refer to [kubebench_design.md](./doc/kubebench_design.md)


## Development

Ensure you run `$ make verify` before submitting PRs. 

// TODO post detailed development guide.
