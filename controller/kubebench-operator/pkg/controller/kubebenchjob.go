package controller

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	argoproj "github.com/argoproj/argo/pkg/client/clientset/versioned/typed/workflow/v1alpha1"

	// workflowUtils "github.com/kubeflow/kubebench/controller/kubebench-operator/util"
	"github.com/kubeflow/kubebench/controller/kubebench-operator/pkg/handler"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type KubebenchJobController struct {
	Logger    *log.Entry
	Clientset kubernetes.Interface
	Queue     workqueue.RateLimitingInterface
	Informer  cache.SharedIndexInformer
	Handler   handler.KubebenchJobHandler
	Workflows argoproj.WorkflowInterface
}

func (c *KubebenchJobController) Run(stopCh <-chan struct{}) {

	defer utilruntime.HandleCrash()
	defer c.Queue.ShutDown()

	c.Logger.Info("Controller.Run: initiating")

	go c.Informer.Run(stopCh)

	if !cache.WaitForCacheSync(stopCh, c.HasSynced) {
		utilruntime.HandleError(fmt.Errorf("Error syncing cache"))
		return
	}
	c.Logger.Info("Controller.Run: cache sync complete")

	wait.Until(c.runWorker, time.Second, stopCh)
}

func (c *KubebenchJobController) HasSynced() bool {
	return c.Informer.HasSynced()
}

func (c *KubebenchJobController) runWorker() {
	log.Info("Controller.runWorker: starting")

	for c.processNextItem() {
		log.Info("Controller.runWorker: processing next item")
	}

	log.Info("Controller.runWorker: completed")
}

func (c *KubebenchJobController) processNextItem() bool {
	log.Info("Controller.processNextItem: start")

	key, quit := c.Queue.Get()

	if quit {
		return false
	}

	defer c.Queue.Done(key)

	keyRaw := key.(string)

	item, exists, err := c.Informer.GetIndexer().GetByKey(keyRaw)
	if err != nil {
		if c.Queue.NumRequeues(key) < 5 {
			c.Logger.Errorf("Controller.processNextItem: Failed processing item with key %s with error %v, retrying", key, err)
			c.Queue.AddRateLimited(key)
		} else {
			c.Logger.Errorf("Controller.processNextItem: Failed processing item with key %s with error %v, no more retries", key, err)
			c.Queue.Forget(key)
			utilruntime.HandleError(err)
		}
	}

	if !exists {
		c.Logger.Infof("Controller.processNextItem: object deleted detected: %s", keyRaw)
		c.Handler.ObjectDeleted(item)
		c.Queue.Forget(key)
	} else {
		c.Logger.Infof("Controller.processNextItem: object created detected: %s", keyRaw)
		c.Handler.ObjectCreated(item)
		c.Queue.Forget(key)
	}

	return true
}
