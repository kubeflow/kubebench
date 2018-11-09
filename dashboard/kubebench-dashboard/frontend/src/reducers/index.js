import * as ActionTypes from '../actions';
// import { combineReducers } from 'redux';

const initialState = {
    yaml: `apiVersion: kubebench.operator/v1
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
  # secretSpecs: # optional
  #   githubTokenSecret: # optional
  #     secretName: my-github-token-secret
  #     secretKey: my-github-token-secret-key
  #   gcpCredentialsSecret: # optional
  #     secretName: my-gcp-credentials-secret
  #     secretKey: my-gcp-credentials-secret-key
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
              name: kubebench-example-tfcnn-with-monitoring
              package: kubebench-examples
              registry: /kubebench/config/registry/kubebench
        manifestParameters:
          valueFrom:
            path: tf-cnn/tf-cnn-dummy.yaml
        createSuccessCondition: createSuccess # optional
        createFailureCondition: createFailure # optional
        runSuccessCondition: runSuccess # optional
        runFailureCondition: runFailure # optional
        #other optional fields: "manifest" - string of raw manifest
    postJob: # optional
      container: # optional between "container" and "resource"
        name: my-postjob
        image: gcr.io/kubeflow-images-public/kubebench/kubebench-example-tf-cnn-post-processor:3c75b50
  reportSpecs: # optional
    csv: # optional
        - inputPath: result.json
          outputPath: report.csv`,
    loading: false,
    snackOpen: false,
    snackText: '',
    filter: '',
    filterType: {
        "Running": true,
        "Failed": true,
        "Succeeded": true,
    },
    jobsList: [
    ],
    filteredJobsList: [
    ],
    modalOpen: false,
    currentId: null,
    currentName: '',
    currentLinks: [

    ],
    parameters: [
        {
            name: "General section",
            description: "section",
        },
        {
            name: "name",
            value: "kubebench-job",
            description: "Job name",
        },
        {
            name: "namespace",
            value: "default",
            description: "Job namespace"
        },
        {
            name: "serviceAccount",
            value: "default",
            description: "The service account used to run the job",
        },
        {
            name: "Secrets section",
            description: "section",
        },
        {
            name: "githubTokenSecret",
            value: "",
            description: "GitHub token secret",
        },
        {
            name: "githubTokenSecretKey",
            value: "",
            description: "Key of GitHub token secret",
        },
        {
            name: "gcpCredentialsSecret",
            value: "",
            description: "GCP credentials secret",
        },
        {
            name: "gcpCredentialsSecretKey",
            value: "",
            description: "Key of GCP credentials secret",
        },
        {
            name: "Main job section",
            description: "section",
        },
        {
            name: "mainJobKsPrototype",
            value: "kubebench-example-tfcnn-with-monitoring",
            description: "The Ksonnet prototype of the job being benchmarked",
        },
        {
            name: "mainJobKsPackage",
            value: "kubebench-examples",
            description: "The Ksonnet package of the job being benchmarked",
        },
        {
            name: "mainJobKsRegistry",
            value: "/kubebench/config/registry/kubebench",
            description: "The Ksonnet registry of the job being benchmarked",
        },
        {
            name: "mainJobConfig",
            value: "tf-cnn/tf-cnn-dummy.yaml",
            description: "Path to the config of the benchmarked job",
        },
        {
            name: "Volumes section",
            description: "section",
        },
        {
            name: "experimentConfigPvc",
            value: "kubebench-config-pvc",
            description: "Configuration PVC",
        },
        {
            name: "experimentDataPvc",
            value: "",
            description: "Data PVC",
        },
        {
            name: "experimentRecordPvc",
            value: "kubebench-exp-pvc",
            description: "Experiment PVC",
        },
        {
            name: "Post job section",
            description: "section",
        },
        {
            name: "postJobImage",
            value: "gcr.io/kubeflow-images-public/kubebench/kubebench-example-tf-cnn-post-processor:3c75b50",
            description: "Image of post processor",
        },
        {
            name: "postJobArgs",
            value: "",
            description: "Arguments of post processor",
        },
        {
            name: "reportType",
            value: "csv",
            description: "Type of reporter",
        },
        {
            name: "csvReporterInput",
            value: "result.json",
            description: "The input of CSV reporter",
        },
        {
            name: "csvReporterOutput",
            value: "report.csv",
            description: "The output of CSV reporter",
        }
    ]
};

const filterValue = (obj, key) => {
    return obj.findIndex(p => p.name === key)
};

const rootReducer = (state = initialState, action) => {
    switch (action.type) {
        // MODIFY
        case ActionTypes.CHANGE_YAML:
            return {
                ...state,
                yaml: action.yaml,
            };
        
        // DEPLOY WHOLE
        case ActionTypes.DEPLOY_SUBMIT:
            return {
                ...state,
                loading: action.loading,
                // snackOpen: true,
                // snackText: 'Successfully deployed',
            };
        case ActionTypes.DEPLOY_SUCCESS:
            return {
                ...state,
                loading: false,
                snackOpen: true,
                snackText: action.text,
            };
        case ActionTypes.DEPLOY_FAILURE:
            return {
                ...state,
                loading: false,
                snackOpen: true,
                snackText: action.error
            };

        
        // DEPLOY PARAMS
        case ActionTypes.DEPLOY_PARAM_SUBMIT:
            return {
                ...state,
                loading: action.loading,
            };
        case ActionTypes.DEPLOY_PARAM_SUCCESS:
            return {
                ...state,
                loading: false,
                snackOpen: true,
                snackText: action.text,
            };
        case ActionTypes.DEPLOY_PARAM_FAILURE:
            return {
                ...state,
                loading: false,
                snackOpen: true,
                snackText: action.error,
            };
        
        // SNACK
        case ActionTypes.CLOSE_SNACK:
            return {
                ...state,
                snackOpen: false,
            };

        // SELECT_JOB
        case ActionTypes.SELECT_JOB:
            return {
                ...state,
                modalOpen: true,
                currentId: action.id,
                currentName: state.jobsList[action.id].name,
            };
        case ActionTypes.CLOSE_SELECT_JOB:
            return {
                ...state,
                modalOpen: false,
                currentId: action.id,
            };
            
        // MODIFY
        case ActionTypes.CHANGE_PARAMETER:
            let params = state.parameters.slice();
            let index = filterValue(params, action.name);
            params[index].value = action.value;
            return {
                ...state,
                parameters: params,
            };

        // FETCH
        case ActionTypes.FETCH_JOB_REQUEST:
            return {
                ...state,
                loading: action.loading,
            };
        case ActionTypes.FETCH_JOB_SUCCESS:
            return {
                ...state,
                jobsList: action.jobsList,
                filteredJobsList: action.jobsList,
                loading: false,
            };
        case ActionTypes.FETCH_JOB_FAILURE:
            return {
                ...state,
                loading: false,
                snackOpen: true,
                snackText: action.error,
            };

        // FILTER 
        case ActionTypes.FILTER_JOBS:
            const jobs = state.jobsList.slice();
            const newList = jobs.filter(job => job.name.includes(action.filter));

            const avTypes = Object.assign({}, state.filterType);
            var typeKeys = Object.keys(avTypes);

            var avFilters = typeKeys.filter((key) => {
                return avTypes[key]
            });
            const filteredJobs = newList.filter(job => avFilters.includes(job.status));
            return {
                ...state,
                filteredJobsList: filteredJobs,
                filter: action.filter,
            };

        // FILTER TYPE
        case ActionTypes.CHANGE_TYPE:
            const types = Object.assign({}, state.filterType)
            types[action.filter] = action.checked;
            var keys = Object.keys(types);

            var filters = keys.filter((key) => {
                return types[key]
            });
            const jobsList = state.jobsList.slice();
            const filtered = jobsList.filter(job => filters.includes(job.status));
            // types[action.type] = 
            return {
                ...state,
                filterType: types,
                filteredJobsList: filtered,
            };

        // DELETE
        case ActionTypes.DELETE_SUBMIT:
            return {
                ...state,
                loading: action.loading,
            };
        case ActionTypes.DELETE_SUCCESS:
            return {
                ...state,
                loading: false,
                snackOpen: true,
                snackText: action.text,
                modalOpen: false,
            };
        case ActionTypes.DELETE_FAILURE: 
            return {
                ...state, 
                loading: false,
                snackOpen: true,
                snackText: action.error,
            };

        default:
            return state;

    }
};

export default rootReducer;