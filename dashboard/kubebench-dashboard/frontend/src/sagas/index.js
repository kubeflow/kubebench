import { take, put, call, fork, select, all, takeEvery } from 'redux-saga/effects';
import * as actions from '../actions';
import axios from 'axios';

export const submitYaml = function *() {
    while (true) {
        const action = yield take(actions.DEPLOY_SUBMIT);
        try {
            const result = yield call(
                submitYamlGo,
                action.yaml,
            );
            if (result === 200) {
                yield put({
                    type: actions.DEPLOY_SUCCESS,
                    text: "Successfully submitted.",
                });
            }
        } catch (err) {
            yield put({
                type: actions.DEPLOY_FAILURE,
                error: "Whoops, something is wrong..."
            })
        }
    }
};

const submitYamlGo = function* (yaml) {
    try {
        let data = {yaml};
        const result = yield call(
            axios.post,
            '/dashboard/submit_yaml/',
            data,
        );

        return result.status
    }catch (err) {
        yield put({
            type: actions.DEPLOY_FAILURE,
            error: "Whoops, something is wrong...",
        });
    }
}

export const submitParams = function *() {
    while (true) {
        const action = yield take(actions.DEPLOY_PARAM_SUBMIT);
        try {
            const result = yield call(
                submitParamsGo,
                action.parameters
            );
            if (result === 200) {
                yield put({
                    type: actions.DEPLOY_PARAM_SUCCESS,
                    text: "Successfully submitted.",
                });
            }
        } catch (err) {
            yield put({
                type: actions.DEPLOY_PARAM_FAILURE,
                error: "Whoops, something is wrong..."
            })
        }
    }
};

const submitParamsGo = function* (parameters) {
    try {
        let newParameters = {};
        parameters.map((item, i) => {
            newParameters[item.name] = item.value;
        })

        const result = yield call(
            axios.post,
            '/dashboard/submit_params/',
            newParameters,
        );
        
        return result.status
    }catch (err) {
        yield put({
            type: actions.DEPLOY_PARAM_FAILURE,
            error: "Whoops, something is wrong...",
        });
    }
}

export const fetchJobsSaga = function *() {
    while (true) {
        const action = yield take(actions.FETCH_JOB_REQUEST);
        try {
            const result = yield call(
                fetchJobsGo
            );
        
            let jobsList = [];
            for(let i = 0; i < result.Names.length; i++) {
                jobsList.push(
                    {
                        name: result.Names[i],
                        status: result.Status[i],
                    }
                );
            }
            yield put({
                type: actions.FETCH_JOB_SUCCESS,
                jobsList: jobsList,
            });
        } catch (err) {
            return Promise.reject(err.message)
        }
    }
};

const fetchJobsGo = function* () {
    try {
        const result = yield call(
            axios.get,
            '/dashboard/fetch_jobs/'
        );

        return result.data;
    } catch (err) {
        yield put({
            type: actions.FETCH_JOB_FAILURE,
            error: "Whoops, something is wrong...",
        });
    }
}

export const deleteJobSaga = function *() {
    while (true) {
        const action = yield take(actions.DELETE_SUBMIT);
        try {
            const result = yield call(
                deleteJobGo,
                action.name
            );
            if (result.Status === "ok") {
                yield put({
                    type: actions.DELETE_SUCCESS,
                    text: "Successfully deleted",
                });
            } else {
                yield put({
                    type: actions.DELETE_FAILURE,
                    error: "Whoops, something is wrong",
                });
            }

        } catch (err) {
            return Promise.reject(err.message);
        }
    }
}

const deleteJobGo = function* (name) {
    try {
        let data = {name};
        const result = yield call(
            axios.post,
            '/dashboard/delete_job/',
            data,
        );

        return result.data
    }catch (err) {
        yield put({
            type: actions.DELETE_FAILURE,
            error: "Whoops, something is wrong...",
        });
    }
}

export default function* rootSaga() {
    yield all([
        fork(fetchJobsSaga),
        fork(deleteJobSaga),
        fork(submitYaml),
        fork(submitParams),
    ]);
};
