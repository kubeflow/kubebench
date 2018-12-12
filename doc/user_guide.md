# User Guide

## Introduction

Kubebench is a harness for benchmarking ML workloads on Kubernetes. Kubebench enables benchmarks by leveraging Kubeflow job operators, as well as Argo workflows.

![Kubebench overview](/doc/images/kubebench_high_level.png)

### Prerequisites

Kubebench runs on a Kubernetes cluster with an existing deployment of Kubeflow core components and Argo. You may refer to [Kubeflow doc](https://www.kubeflow.org/docs/about/kubeflow/) and [Argo doc](https://argoproj.github.io/docs/argo/readme.html) for details.

### Glossary

- Kubebench Job

  A Kubebench Job is a workflow that runs a benchmark experiment. The Kubebench Job consists of several steps, that include config loading, Kubeflow job generation, benchmark running, result collection and aggregation, etc.

- Kubeflow Job

  A Kubeflow Job is the workload being benchmarked. The Kubeflow Job can be a TFJob or a PyTorchJob (more Kubeflow specific job types will be supported soon). The Kubeflow job is run as one step in the Kubebench workflow.

- Experiment

  An Experiment is one single run of benchmark for a Kubeflow Job. When you run an experiment, both the workflow (Kubebench Job) and the workload (Kubeflow Job) will generate lots of data and information associated with the particular experiment. Kubebench will persist these info automatically in a per experiment basis.

- Job Template & Parameters

  A Kubebench Job in a benchmark experiment can be uniquely defined by a template and a set of parameters. The template generates a manifest file (a Kubernetes resource descripter) with configurable variables, and the parameters provide specific values of the variables. Currently, Kubebench supports Ksonnet prototype as the template format, and a YAML config as the parameters.

- Report

  A Report is an aggregated list of results from multiple Experiments. The report can be in the form of a database, or simply a local file. Currently Kubebench support CSV file based report.

## Preparing Storage

In order to run benchmarks with Kubebench, you need to make your data and configs available to Kubebench by preparing them in Persistent Volumes (PV). Kubebench expects the following user provided volumes:

- Config volume (required): stores all the parameters and (optionally) job templates. You need to store your experiment configurations in this volume before running Kubebench job.

- Experiment volume (required): stores all the detailed information of each experiment during runtime, including configs, intermediate outputs, and final results. You will just provide an empty volume to Kubebench as the experiment volume, and Kubebench will automatically populate it during each experiments.

- Data volume (optional): if your benchmark job needs to access static data, then you can put your data in this volume and specify it in Kubebench job config. Kubebench will automatically mount the volume and make the data available to your job.

You may find an example Kubebench directory structure in the Appendices.

(Note: When you install the quick starter package, it will automatically prepare the storage for you with a couple of example job configs in an NFS container.)

## Writing Benchmark Codes

The benchmark codes live in 2 components: the main-job and the post-job. The main-job supports either a Kubeflow Job, or a native Kubernetes Job. The post-job supports a native Kubernetes Job. You need to implement codes in these 2 components and provide them to Kubebench as Docker images.

![Writing Kubebench jobs](/doc/images/writing_kubebench_jobs.png)

The diagram above shows a rough idea of how your main-job and post-job will interact with each other and with the rest of the Kubebench workflow. When your job is deployed through Kubebench, all the containers of your job will automatically mount the volumes you configured. They will also have the following environment variables available so as to make it easy to get experiment data and share information between jobs.

Name | Description | Default Value
--- | --- | ---
KUBEBENCH_CONFIG_ROOT | The root path of all job configs | N/A
KUBEBENCH_DATA_ROOT | The root path of all data | N/A
KUBEBENCH_EXP_ROOT | The root path of all experiments | N/A
KUBEBENCH_EXP_ID | The ID of a particular experiment | N/A
KUBEBENCH_EXP_PATH | The root path of a particular experiment | $KUBEBENCH_EXP_ROOT/$KUBEBENCH_EXP_ID
KUBEBENCH_EXP_CONFIG_PATH | The path of a particular experiment's config | $KUBEBENCH_EXP_PATH/config
KUBEBENCH_EXP_OUTPUT_PATH | The path of a particular experiment's job outputs | $KUBEBENCH_EXP_PATH/output
KUBEBENCH_EXP_RESULT_PATH | The path of a particular experiment's result | $KUBEBENCH_EXP_PATH/result

When writing codes for the main-job and post-job, please follow a few basic rules, so that the Kubebench workflow can function properly:

The main job need to:

- run the benchmark codes
- write outputs to `${KUBEBENCH_EXP_OUTPUT_PATH}`

The post job need to:

- read main job outputs from `${KUBEBENCH_EXP_OUTPUT_PATH}`
- parse the outputs and construct a json formated result file with desired information
- write the result to `${KUBEBENCH_EXP_RESULT_PATH}`

Once the result file is available, the Kubebench reporter will automatically pick it up and report the results to user specified destinations.

## Configuring Benchmark Jobs

The Kubebench jobs can be configured through Ksonnet. To create a Ksonnet component, do the followings in your Ksonnet app with an existing Kubeflow installation.

```bash
ks pkg install kubeflow/kubebench
ks generate kubebench-job <JOB_NAME>
```

You can then set each parameter of your Kubebench job in the following way.

```bash
ks param set <JOB_NAME> <PARAM_KEY> <PARAM_VALUE> --env=<KS_ENV>
```

(Note: replace `<JOB_NAME>`, `<PARAM_KEY>`, `<PARAM_VALUE>`, `<KS_ENV>` with your own values)

 Please refer to [Kubeflow doc](https://www.kubeflow.org/docs/guides/components/ksonnet/) for further details about how to use Ksonnet.

### Set volume parameters

Once the volumes are prepared, create a Persistent Volume Claim (PVC) for each volume and give the PVC names in the following parameters in your Ksonnet component config.

- `experimentConfigPvc`: name of the PVC pointing to your config volume
- `experimentRecordPvc`: name of the PVC pointing to your experiment volume
- `experimentDataPvc`: name of the PVC pointing to your data volume

### Set job parameters

#### Main job

The main job requires a unique reference to a Ksonnet prototype (i.e. registry, package, and prototype name) and a path to the parameter config file:

- `mainJobKsRegistry`: location of main job's Ksonnet registry
- `mainJobKsPackage`: main job's Ksonnet package
- `mainJobKsPrototype`: main job's Ksonnet prototype
- `mainJobConfig`: main job's parameters

The parameter config file should be located in your config volume and the path given should be relative to the config volume's mount point. If using a file path as the Ksonnet registry, the path given should be relative to the config volume's mount point. If using a github repository as the Ksonnet registry, you may need to set a github secret to avoid hitting API quota limit, and provide the following parameter values to Kubebench.

- `githubTokenSecret`: the name of github token secret
- `githubTokenSecretKey`: the key of the secret to retrieve github token value

#### Post job

The post job is deployed as a native Kubernetes job. You can specify the image and arguements used in the job.

- `postJobImage`: the image of the post job
- `postJobArgs`: the arguments of the post job

### Set reporter parameters

When you run multiple benchmark experiments, Kubebench reporter can aggregate your experiment results into a single dataset. Currently Kubebench supports result aggregation into a CSV formated file stored in your experiment volume. You may specify the following reporter parameters to configure the reporter.

- `csvReporterInput`: the input of the csv reporter (i.e. the output file of post job)
- `csvReporterOutput`: the output of the csv reporter

Note that the `csvReporterInput` is a path relative to `$KUBEBENCH_EXP_RESULT_PATH`, and the `csvReporterOutput` is a path relative to `$KUBEBENCH_EXP_ROOT`.

## Running Kubebench Jobs

### Start a job

Once you have configured the parameters of your Kubebench job, you can start it with

```bash
ks apply <KS_ENV> -c <JOB_NAME>
```

### Check job status and results

The Kubebench job is deployed as an Argo workflow. When the job is running, you can go to the Argo UI to keep track of the job progress. Once the job is finished, you will find an experiment specific directory in the experiment volume, where you will find all the information related with the particular experiment.

### Clean up a job

You may delete the Kubebench job with

```bash
ks delete <KS_ENV> -c <JOB_NAME>
```

## Appendices

### Example job templates and parameters

The example job templates (Ksonnet prototype) can be found [here](/kubebench/kubebench-examples).

The example parameters (YAML file) can be found [here](/examples/config).

### Example Kubebench directory structure

The followings show an example Kubebench directory.

- The `config` and `data` are prepared by user.
- If you want to use local Ksonnet registry, you can place the registry in an arbitrary subdirectory under `config`, and provide the relative path to Kubebench job.
- The `experiments` is automatically populated by Kubebench, where each experiment will have a unique ID and its data will be kept in a dedicated subdirectory.

```
/kubebench
├── config
│   └── tf-cnn-dummy.yaml
├── data
│   └── train_data.tfrecords
└── experiments
    ├── report.csv
    └── tf-cnn-dummy-201809150923-z72k
        ├── config
        │   ├── kf-job-manifest.yaml
        │   └── tf-cnn-dummy.yaml
        ├── output
        │   ├── worker0.log
        │   ├── worker1.log
        │   └── ps0.log
        └── result
            └── result.json
```
