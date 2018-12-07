// Copyright 2018 The Kubeflow Authors
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

// Package controller provides a Kubernetes controller for a TFJob resource.
package tensorflow

import (
	"testing"

	"k8s.io/api/core/v1"

	tfv1alpha2 "github.com/kubeflow/tf-operator/pkg/apis/tensorflow/v1alpha2"
	"github.com/kubeflow/tf-operator/pkg/util/testutil"
)

func TestFailed(t *testing.T) {
	tfJob := testutil.NewTFJob(3, 0)
	initializeTFReplicaStatuses(tfJob, tfv1alpha2.TFReplicaTypeWorker)
	pod := testutil.NewBasePod("pod", tfJob, t)
	pod.Status.Phase = v1.PodFailed
	updateTFJobReplicaStatuses(tfJob, tfv1alpha2.TFReplicaTypeWorker, pod)
	if tfJob.Status.TFReplicaStatuses[tfv1alpha2.TFReplicaTypeWorker].Failed != 1 {
		t.Errorf("Failed to set the failed to 1")
	}
	err := updateStatusSingle(tfJob, tfv1alpha2.TFReplicaTypeWorker, 3, false, false)
	if err != nil {
		t.Errorf("Expected error %v to be nil", err)
	}
	found := false
	for _, condition := range tfJob.Status.Conditions {
		if condition.Type == tfv1alpha2.TFJobFailed {
			found = true
		}
	}
	if !found {
		t.Errorf("Failed condition is not found")
	}
}

func TestStatus(t *testing.T) {
	type testCase struct {
		description string
		tfJob       *tfv1alpha2.TFJob

		expectedFailedPS    int32
		expectedSucceededPS int32
		expectedActivePS    int32

		expectedFailedWorker    int32
		expectedSucceededWorker int32
		expectedActiveWorker    int32

		expectedFailedChief    int32
		expectedSucceededChief int32
		expectedActiveChief    int32

		restart          bool
		worker0Completed bool

		expectedType tfv1alpha2.TFJobConditionType
	}

	testCases := []testCase{
		testCase{
			description:             "Chief worker is succeeded",
			tfJob:                   testutil.NewTFJobWithChief(1, 0),
			expectedFailedPS:        0,
			expectedSucceededPS:     0,
			expectedActivePS:        0,
			expectedFailedWorker:    0,
			expectedSucceededWorker: 1,
			expectedActiveWorker:    0,
			expectedFailedChief:     0,
			expectedSucceededChief:  1,
			expectedActiveChief:     0,
			restart:                 false,
			worker0Completed:        false,
			expectedType:            tfv1alpha2.TFJobSucceeded,
		},
		testCase{
			description:             "Chief worker is running",
			tfJob:                   testutil.NewTFJobWithChief(1, 0),
			expectedFailedPS:        0,
			expectedSucceededPS:     0,
			expectedActivePS:        0,
			expectedFailedWorker:    0,
			expectedSucceededWorker: 0,
			expectedActiveWorker:    0,
			expectedFailedChief:     0,
			expectedSucceededChief:  0,
			expectedActiveChief:     1,
			restart:                 false,
			worker0Completed:        false,
			expectedType:            tfv1alpha2.TFJobRunning,
		},
		testCase{
			description:             "Chief worker is failed",
			tfJob:                   testutil.NewTFJobWithChief(1, 0),
			expectedFailedPS:        0,
			expectedSucceededPS:     0,
			expectedActivePS:        0,
			expectedFailedWorker:    0,
			expectedSucceededWorker: 0,
			expectedActiveWorker:    0,
			expectedFailedChief:     1,
			expectedSucceededChief:  0,
			expectedActiveChief:     0,
			restart:                 false,
			worker0Completed:        false,
			expectedType:            tfv1alpha2.TFJobFailed,
		},
		testCase{
			description:             "(No chief worker) Worker is failed",
			tfJob:                   testutil.NewTFJob(1, 0),
			expectedFailedPS:        0,
			expectedSucceededPS:     0,
			expectedActivePS:        0,
			expectedFailedWorker:    1,
			expectedSucceededWorker: 0,
			expectedActiveWorker:    0,
			expectedFailedChief:     0,
			expectedSucceededChief:  0,
			expectedActiveChief:     0,
			restart:                 false,
			worker0Completed:        false,
			expectedType:            tfv1alpha2.TFJobFailed,
		},
		testCase{
			description:             "(No chief worker) Worker is succeeded",
			tfJob:                   testutil.NewTFJob(1, 0),
			expectedFailedPS:        0,
			expectedSucceededPS:     0,
			expectedActivePS:        0,
			expectedFailedWorker:    0,
			expectedSucceededWorker: 1,
			expectedActiveWorker:    0,
			expectedFailedChief:     0,
			expectedSucceededChief:  0,
			expectedActiveChief:     0,
			restart:                 false,
			worker0Completed:        false,
			expectedType:            tfv1alpha2.TFJobSucceeded,
		},
		testCase{
			description:             "(No chief worker) Worker is running",
			tfJob:                   testutil.NewTFJob(1, 0),
			expectedFailedPS:        0,
			expectedSucceededPS:     0,
			expectedActivePS:        0,
			expectedFailedWorker:    0,
			expectedSucceededWorker: 0,
			expectedActiveWorker:    1,
			expectedFailedChief:     0,
			expectedSucceededChief:  0,
			expectedActiveChief:     0,
			restart:                 false,
			worker0Completed:        false,
			expectedType:            tfv1alpha2.TFJobRunning,
		},
		testCase{
			description:             "(No chief worker) 2 workers are succeeded, 2 workers are active",
			tfJob:                   testutil.NewTFJob(4, 2),
			expectedFailedPS:        0,
			expectedSucceededPS:     0,
			expectedActivePS:        2,
			expectedFailedWorker:    0,
			expectedSucceededWorker: 2,
			expectedActiveWorker:    2,
			expectedFailedChief:     0,
			expectedSucceededChief:  0,
			expectedActiveChief:     0,
			restart:                 false,
			worker0Completed:        false,
			expectedType:            tfv1alpha2.TFJobRunning,
		},
		testCase{
			description:             "(No chief worker) 2 workers are running, 2 workers are failed",
			tfJob:                   testutil.NewTFJob(4, 2),
			expectedFailedPS:        0,
			expectedSucceededPS:     0,
			expectedActivePS:        2,
			expectedFailedWorker:    2,
			expectedSucceededWorker: 0,
			expectedActiveWorker:    2,
			expectedFailedChief:     0,
			expectedSucceededChief:  0,
			expectedActiveChief:     0,
			restart:                 false,
			worker0Completed:        false,
			expectedType:            tfv1alpha2.TFJobFailed,
		},
		testCase{
			description:             "(No chief worker) 2 workers are succeeded, 2 workers are failed",
			tfJob:                   testutil.NewTFJob(4, 2),
			expectedFailedPS:        0,
			expectedSucceededPS:     0,
			expectedActivePS:        2,
			expectedFailedWorker:    2,
			expectedSucceededWorker: 2,
			expectedActiveWorker:    0,
			expectedFailedChief:     0,
			expectedSucceededChief:  0,
			expectedActiveChief:     0,
			restart:                 false,
			worker0Completed:        false,
			expectedType:            tfv1alpha2.TFJobFailed,
		},
		testCase{
			description:             "(No chief worker) worker-0 are succeeded, 3 workers are active",
			tfJob:                   testutil.NewTFJob(4, 2),
			expectedFailedPS:        0,
			expectedSucceededPS:     0,
			expectedActivePS:        2,
			expectedFailedWorker:    0,
			expectedSucceededWorker: 1,
			expectedActiveWorker:    3,
			expectedFailedChief:     0,
			expectedSucceededChief:  0,
			expectedActiveChief:     0,
			restart:                 false,
			worker0Completed:        true,
			expectedType:            tfv1alpha2.TFJobSucceeded,
		},
		testCase{
			description:             "Chief is running, workers are failed",
			tfJob:                   testutil.NewTFJobWithChief(4, 2),
			expectedFailedPS:        0,
			expectedSucceededPS:     0,
			expectedActivePS:        2,
			expectedFailedWorker:    4,
			expectedSucceededWorker: 0,
			expectedActiveWorker:    0,
			expectedFailedChief:     0,
			expectedSucceededChief:  0,
			expectedActiveChief:     1,
			restart:                 false,
			worker0Completed:        false,
			expectedType:            tfv1alpha2.TFJobRunning,
		},
		testCase{
			description:             "Chief is running, workers are succeeded",
			tfJob:                   testutil.NewTFJobWithChief(4, 2),
			expectedFailedPS:        0,
			expectedSucceededPS:     0,
			expectedActivePS:        2,
			expectedFailedWorker:    0,
			expectedSucceededWorker: 4,
			expectedActiveWorker:    0,
			expectedFailedChief:     0,
			expectedSucceededChief:  0,
			expectedActiveChief:     1,
			restart:                 false,
			worker0Completed:        false,
			expectedType:            tfv1alpha2.TFJobRunning,
		},
		testCase{
			description:             "Chief is running, a PS is failed",
			tfJob:                   testutil.NewTFJobWithChief(4, 2),
			expectedFailedPS:        1,
			expectedSucceededPS:     0,
			expectedActivePS:        1,
			expectedFailedWorker:    0,
			expectedSucceededWorker: 4,
			expectedActiveWorker:    0,
			expectedFailedChief:     0,
			expectedSucceededChief:  0,
			expectedActiveChief:     1,
			restart:                 false,
			worker0Completed:        false,
			expectedType:            tfv1alpha2.TFJobFailed,
		},
		testCase{
			description:             "Chief is failed, workers are succeeded",
			tfJob:                   testutil.NewTFJobWithChief(4, 2),
			expectedFailedPS:        0,
			expectedSucceededPS:     0,
			expectedActivePS:        2,
			expectedFailedWorker:    0,
			expectedSucceededWorker: 4,
			expectedActiveWorker:    0,
			expectedFailedChief:     1,
			expectedSucceededChief:  0,
			expectedActiveChief:     0,
			restart:                 false,
			worker0Completed:        false,
			expectedType:            tfv1alpha2.TFJobFailed,
		},
		testCase{
			description:             "Chief is succeeded, workers are failed",
			tfJob:                   testutil.NewTFJobWithChief(4, 2),
			expectedFailedPS:        0,
			expectedSucceededPS:     0,
			expectedActivePS:        2,
			expectedFailedWorker:    4,
			expectedSucceededWorker: 0,
			expectedActiveWorker:    0,
			expectedFailedChief:     0,
			expectedSucceededChief:  1,
			expectedActiveChief:     0,
			restart:                 false,
			worker0Completed:        false,
			expectedType:            tfv1alpha2.TFJobSucceeded,
		},
		testCase{
			description:             "Chief is failed and restarting",
			tfJob:                   testutil.NewTFJobWithChief(4, 2),
			expectedFailedPS:        0,
			expectedSucceededPS:     0,
			expectedActivePS:        2,
			expectedFailedWorker:    4,
			expectedSucceededWorker: 0,
			expectedActiveWorker:    0,
			expectedFailedChief:     1,
			expectedSucceededChief:  0,
			expectedActiveChief:     0,
			restart:                 true,
			worker0Completed:        false,
			expectedType:            tfv1alpha2.TFJobRestarting,
		},
	}

	for i, c := range testCases {
		initializeTFReplicaStatuses(c.tfJob, tfv1alpha2.TFReplicaTypeWorker)
		initializeTFReplicaStatuses(c.tfJob, tfv1alpha2.TFReplicaTypeChief)
		initializeTFReplicaStatuses(c.tfJob, tfv1alpha2.TFReplicaTypePS)

		setStatusForTest(c.tfJob, tfv1alpha2.TFReplicaTypePS, c.expectedFailedPS, c.expectedSucceededPS, c.expectedActivePS, t)
		setStatusForTest(c.tfJob, tfv1alpha2.TFReplicaTypeWorker, c.expectedFailedWorker, c.expectedSucceededWorker, c.expectedActiveWorker, t)
		setStatusForTest(c.tfJob, tfv1alpha2.TFReplicaTypeChief, c.expectedFailedChief, c.expectedSucceededChief, c.expectedActiveChief, t)

		if _, ok := c.tfJob.Spec.TFReplicaSpecs[tfv1alpha2.TFReplicaTypeChief]; ok {
			err := updateStatusSingle(c.tfJob, tfv1alpha2.TFReplicaTypeChief, 1, c.restart, c.worker0Completed)
			if err != nil {
				t.Errorf("%s: Expected error %v to be nil", c.description, err)
			}
			if c.tfJob.Spec.TFReplicaSpecs[tfv1alpha2.TFReplicaTypeWorker] != nil {
				replicas := c.tfJob.Spec.TFReplicaSpecs[tfv1alpha2.TFReplicaTypeWorker].Replicas
				err := updateStatusSingle(c.tfJob, tfv1alpha2.TFReplicaTypeWorker, int(*replicas), c.restart, c.worker0Completed)
				if err != nil {
					t.Errorf("%s: Expected error %v to be nil", c.description, err)
				}
			}
			if c.tfJob.Spec.TFReplicaSpecs[tfv1alpha2.TFReplicaTypePS] != nil {
				replicas := c.tfJob.Spec.TFReplicaSpecs[tfv1alpha2.TFReplicaTypePS].Replicas
				err := updateStatusSingle(c.tfJob, tfv1alpha2.TFReplicaTypePS, int(*replicas), c.restart, c.worker0Completed)
				if err != nil {
					t.Errorf("%s: Expected error %v to be nil", c.description, err)
				}
			}
		} else {
			if c.tfJob.Spec.TFReplicaSpecs[tfv1alpha2.TFReplicaTypeWorker] != nil {
				replicas := c.tfJob.Spec.TFReplicaSpecs[tfv1alpha2.TFReplicaTypeWorker].Replicas
				err := updateStatusSingle(c.tfJob, tfv1alpha2.TFReplicaTypeWorker, int(*replicas), c.restart, c.worker0Completed)
				if err != nil {
					t.Errorf("%s: Expected error %v to be nil", c.description, err)
				}
			}
			if c.tfJob.Spec.TFReplicaSpecs[tfv1alpha2.TFReplicaTypePS] != nil {
				replicas := c.tfJob.Spec.TFReplicaSpecs[tfv1alpha2.TFReplicaTypePS].Replicas
				err := updateStatusSingle(c.tfJob, tfv1alpha2.TFReplicaTypePS, int(*replicas), c.restart, c.worker0Completed)
				if err != nil {
					t.Errorf("%s: Expected error %v to be nil", c.description, err)
				}
			}
		}

		// Test filterOutCondition
		filterOutConditionTest(c.tfJob.Status, t)

		found := false
		for _, condition := range c.tfJob.Status.Conditions {
			if condition.Type == c.expectedType {
				found = true
			}
		}
		if !found {
			t.Errorf("Case[%d]%s: Condition %s is not found", i, c.description, c.expectedType)
		}
	}
}

func setStatusForTest(tfJob *tfv1alpha2.TFJob, typ tfv1alpha2.TFReplicaType, failed, succeeded, active int32, t *testing.T) {
	pod := testutil.NewBasePod("pod", tfJob, t)
	var i int32
	for i = 0; i < failed; i++ {
		pod.Status.Phase = v1.PodFailed
		updateTFJobReplicaStatuses(tfJob, typ, pod)
	}
	for i = 0; i < succeeded; i++ {
		pod.Status.Phase = v1.PodSucceeded
		updateTFJobReplicaStatuses(tfJob, typ, pod)
	}
	for i = 0; i < active; i++ {
		pod.Status.Phase = v1.PodRunning
		updateTFJobReplicaStatuses(tfJob, typ, pod)
	}
}

func filterOutConditionTest(status tfv1alpha2.TFJobStatus, t *testing.T) {
	flag := isFailed(status) || isSucceeded(status)
	for _, condition := range status.Conditions {
		if flag && condition.Type == tfv1alpha2.TFJobRunning && condition.Status == v1.ConditionTrue {
			t.Error("Error condition status when succeeded or failed")
		}
	}
}
