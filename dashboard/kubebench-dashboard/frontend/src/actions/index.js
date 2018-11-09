export const CHANGE_YAML = "CHANGE_YAML";

export const CLOSE_SNACK = "CLOSE_SNACK";

export const SELECT_JOB = "SELECT_JOB";
export const CLOSE_SELECT_JOB = "CLOSE_SELECT_JOB";

export const DEPLOY_SUBMIT = "DEPLOY_SUBMIT";
export const DEPLOY_SUCCESS = "DEPLOY_SUCCESS";
export const DEPLOY_FAILURE = "DEPLOY_FAILURE";

export const DEPLOY_PARAM_SUBMIT = "DEPLOY_PARAM_SUBMIT";
export const DEPLOY_PARAM_SUCCESS = "DEPLOY_PARAM_SUCCESS";
export const DEPLOY_PARAM_FAILURE = "DEPLOY_PARAM_FAILURE";

export const DELETE_SUBMIT = "DELETE SUBMIT";
export const DELETE_SUCCESS = "DELETE_SUCCESS";
export const DELETE_FAILURE = "DELETE_FAILURE";

export const CHANGE_PARAMETER = "CHANGE_PARAMETER";

export const FETCH_JOB_REQUEST = "FETCH_JOB_REQUEST";
export const FETCH_JOB_SUCCESS = "FETCH_JOB_SUCCESS";
export const FETCH_JOB_FAILURE = "FETCH_JOB_FAILURE";

export const FILTER_JOBS = "FILTER_JOBS";
export const CHANGE_TYPE = "CHANGE_TYPE";


export const changeYaml = (yaml) => {
    return {
        type: CHANGE_YAML,
        yaml: yaml,
    };
};

export const changeParameters = (name, value) => {
    return {
        type: CHANGE_PARAMETER,
        name,
        value,
    };
};

export const submitWholeYaml = (yaml) => {
    return {
        type: DEPLOY_SUBMIT,
        loading: true,
        yaml: yaml
    };
};

export const submitYamlParameters = (params) => {
    return {
        type: DEPLOY_PARAM_SUBMIT,
        loading: true,
        parameters: params,
    };
};

export const closeSnack =() => {
    return {
        type: CLOSE_SNACK,
    };
};

export const openModal = (id) => {
    return {
        type: SELECT_JOB,
        id,
    };
};

export const closeModal = () => {
    return {
        type: CLOSE_SELECT_JOB,
    };
};

export const deleteJob = (name) => {
    return {
        type: DELETE_SUBMIT,
        name,
    };
};

export const fetchJobs = () => {
    return {
        type: FETCH_JOB_REQUEST,
        loading: true,
    }
}

export const filterJobs = (filter) => {
    return {
        type: FILTER_JOBS,
        filter,
    };
};

export const changeType = (filter, checked) => {
    return {
        type: CHANGE_TYPE,
        filter,
        checked,
    };
};