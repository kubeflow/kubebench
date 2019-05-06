// Copyright 2019 The Kubeflow Authors
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

package v1alpha2

import (
	"fmt"
	"time"

	argoclientset "github.com/argoproj/argo/pkg/client/clientset/versioned"
	argoinformer "github.com/argoproj/argo/pkg/client/informers/externalversions/workflow/v1alpha1"
	argolister "github.com/argoproj/argo/pkg/client/listers/workflow/v1alpha1"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"

	kbjobv1alpha2 "github.com/kubeflow/kubebench/controller/pkg/apis/kubebenchjob/v1alpha2"
	kbjobclientset "github.com/kubeflow/kubebench/controller/pkg/client/clientset/versioned"
	kbjobscheme "github.com/kubeflow/kubebench/controller/pkg/client/clientset/versioned/scheme"
	kbjobinformer "github.com/kubeflow/kubebench/controller/pkg/client/informers/externalversions/kubebenchjob/v1alpha2"
	kbjoblister "github.com/kubeflow/kubebench/controller/pkg/client/listers/kubebenchjob/v1alpha2"
	kbworkflow "github.com/kubeflow/kubebench/controller/pkg/workflow"
)

const (
	// EventResourceCreated is an Event "reason" when a resource is created
	EventResourceCreated = "Created"
	// EventResourceNotCreated is an Event "reason" when a resource fails to be created
	EventResourceNotCreated = "ErrResourceNotCreated"
	// EventResourceSynced is an Event "reason" when a resource is synced
	EventResourceSynced = "Synced"
	// EventResourceExists is an Event "reason" when a resource fails to sync
	EventResourceExists = "ErrResourceExists"

	// EventMsgWorkflowCreated is an Event "message" when an Argo Workflow owned
	// by a KubebenchJob is created
	EventMsgWorkflowCreated = "Workflow created successfully"
	// EventMsgWorkflowNotCreated is an Event "message" when an Argo Workflow owned
	// by a KubebenchJob fails to be created
	EventMsgWorkflowNotCreated = "Workflow failed to be created"
	// EventMsgWorkflowExists is an Event "message" when a Kubebenchjob fails to sync
	// due to an Argo Workflow already existing
	EventMsgWorkflowExists = "Workflow already exists and is not managed by KubebenchJob"
	// EventMsgKubebenchJobSynced is the message used for an Event fired when a KubebenchJob
	// is synced successfully
	EventMsgKubebenchJobSynced = "KubebenchJob synced successfully"
)

// KubebenchJobController is the controller for KubebenchJob
type KubebenchJobController struct {
	argoWorkflowClientSet      argoclientset.Interface
	kubebenchJobClientSet      kbjobclientset.Interface
	argoWorkflowLister         argolister.WorkflowLister
	kubebenchJobLister         kbjoblister.KubebenchJobLister
	argoWorkflowInformerSynced cache.InformerSynced
	kubebenchJobInformerSynced cache.InformerSynced
	queue                      workqueue.RateLimitingInterface
	recorder                   record.EventRecorder
	config                     kbjobv1alpha2.KubebenchConfig
}

// NewKubebenchJobController creates a new KubebenchJobController
func NewKubebenchJobController(
	argoClientSet argoclientset.Interface,
	kbJobClientSet kbjobclientset.Interface,
	argoInformer argoinformer.WorkflowInformer,
	kbJobInformer kbjobinformer.KubebenchJobInformer,
	config kbjobv1alpha2.KubebenchConfig) *KubebenchJobController {

	utilruntime.Must(kbjobscheme.AddToScheme(scheme.Scheme))
	kbJobController := &KubebenchJobController{
		argoWorkflowClientSet:      argoClientSet,
		kubebenchJobClientSet:      kbJobClientSet,
		argoWorkflowLister:         argoInformer.Lister(),
		kubebenchJobLister:         kbJobInformer.Lister(),
		argoWorkflowInformerSynced: argoInformer.Informer().HasSynced,
		kubebenchJobInformerSynced: kbJobInformer.Informer().HasSynced,
		queue: workqueue.NewNamedRateLimitingQueue(
			workqueue.DefaultControllerRateLimiter(), "kubebenchjob"),
		recorder: record.NewBroadcaster().NewRecorder(
			scheme.Scheme, corev1.EventSource{Component: "kubebenchjob-controller"}),
	}

	kbJobInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: kbJobController.enqueueKubebenchJob,
		UpdateFunc: func(old, new interface{}) {
			kbJobController.enqueueKubebenchJob(new)
		},
	})

	return kbJobController
}

// Run starts a KubebenchJobController
func (c *KubebenchJobController) Run(threadiness int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	if ok := cache.WaitForCacheSync(stopCh, c.argoWorkflowInformerSynced,
		c.kubebenchJobInformerSynced); !ok {
		return fmt.Errorf("Failed to wait for caches to sync")
	}
	log.Info("Controller.Run: cache sync complete.")

	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	log.Info("Started controller workers.")
	<-stopCh
	log.Info("Shutting down controller workers.")

	return nil
}

func (c *KubebenchJobController) runWorker() {
	for c.processNextItem() {
	}
}

func (c *KubebenchJobController) processNextItem() bool {
	obj, quit := c.queue.Get()
	if quit {
		return false
	}

	err := func(obj interface{}) error {
		defer c.queue.Done(obj)
		var key string
		var ok bool
		if key, ok = obj.(string); !ok {
			c.queue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}

		if err := c.syncKubebenchJob(key); err != nil {
			c.queue.AddRateLimited(key)
			return fmt.Errorf("error syncing '%s': %s, requeuing", key, err.Error())
		}

		c.queue.Forget(obj)
		return nil
	}(obj)
	if err != nil {
		utilruntime.HandleError(err)
		return true
	}

	return true
}

func (c *KubebenchJobController) enqueueKubebenchJob(kbjob interface{}) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(kbjob)
	if err != nil {
		utilruntime.HandleError(
			fmt.Errorf("Failed to get key for KubebenchJob object %#v: %s", kbjob, err))
		return
	}
	c.queue.Add(key)
}

func (c *KubebenchJobController) syncKubebenchJob(key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}
	if namespace == "" || name == "" {
		return fmt.Errorf("Invalid key %q", key)
	}

	// Get the kubebenchjob
	kbjob, err := c.kubebenchJobLister.KubebenchJobs(namespace).Get(name)
	if err != nil {
		// Stop processing if the resource no longer exist
		if errors.IsNotFound(err) {
			utilruntime.HandleError(fmt.Errorf("kubebenchjob %q in work queue no longer exists", key))
			return nil
		}
		return err
	}

	// Get info of the workflow owned by the kubebenchjob.
	// The workflow should have the same name and namespace as the kubebenchjob
	workflow, err := c.argoWorkflowLister.Workflows(namespace).Get(name)
	// If workflow does not exist then create
	if errors.IsNotFound(err) {
		wf, err := kbworkflow.BuildWorkflow(kbjob, &c.config, true)
		if err != nil {
			log.Errorf("Failed to build Workflow: %s", err)
			c.recorder.Event(kbjob, corev1.EventTypeWarning, EventResourceNotCreated, EventMsgWorkflowNotCreated)
			return err
		}
		workflow, err = c.argoWorkflowClientSet.ArgoprojV1alpha1().Workflows(namespace).Create(wf)
		if err != nil {
			c.recorder.Event(kbjob, corev1.EventTypeWarning, EventResourceNotCreated, EventMsgWorkflowNotCreated)
			return err
		}
		log.Infof("Created Workflow for KubebenchJob %s", kbjob.Name)
		c.recorder.Event(kbjob, corev1.EventTypeNormal, EventResourceCreated, EventMsgWorkflowCreated)
	}

	// If the workflow already exists and is not owned by the kubebenchjob, record a warning
	if !metav1.IsControlledBy(workflow, kbjob) {
		c.recorder.Event(kbjob, corev1.EventTypeWarning, EventResourceExists, EventMsgWorkflowExists)
		return fmt.Errorf(EventMsgWorkflowExists)
	}

	// Update the kubebenchjob status to be the same as workflow status
	kbjobCopy := kbjob.DeepCopy()
	kbjob.Status = kbjobv1alpha2.KubebenchJobStatus{WorkflowStatus: workflow.Status}
	_, err = c.kubebenchJobClientSet.KubeflowV1alpha2().KubebenchJobs(namespace).Update(kbjobCopy)
	if err != nil {
		return err
	}

	c.recorder.Event(kbjob, corev1.EventTypeNormal, EventResourceSynced, EventMsgKubebenchJobSynced)
	return nil
}
