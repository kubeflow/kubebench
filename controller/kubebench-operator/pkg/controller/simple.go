package controller

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/kubeflow/kubebench/controller/kubebench-operator/pkg/handler"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

// Controller struct defines how a controller should encapsulate
// logging, client connectivity, informing (list and watching)
// Queueing, and handling of resource changes
type Controller struct {
	Logger    *log.Entry
	Clientset kubernetes.Interface
	Queue     workqueue.RateLimitingInterface
	Informer  cache.SharedIndexInformer
	Handler   handler.SimpleHandler
}

func (c *Controller) Run(stopCh <-chan struct{}) {
	// handle a panic with logging and exiting
	defer utilruntime.HandleCrash()
	// ignore new items in the Queue but when all goroutines
	// have completed existing items then shutdown
	defer c.Queue.ShutDown()

	c.Logger.Info("Controller.Run: initiating")

	// run the informer to start listing and watching resources
	go c.Informer.Run(stopCh)

	// do the initial synchronization (one time) to populate resources
	if !cache.WaitForCacheSync(stopCh, c.HasSynced) {
		utilruntime.HandleError(fmt.Errorf("Error syncing cache"))
		return
	}
	c.Logger.Info("Controller.Run: cache sync complete")

	// run the runWorker method every second with a stop channel
	wait.Until(c.runWorker, time.Second, stopCh)
}

// HasSynced allows us to satisfy the Controller interface
// by wiring up the informer's HasSynced method to it
func (c *Controller) HasSynced() bool {
	return c.Informer.HasSynced()
}

// runWorker executes the loop to process new items added to the Queue
func (c *Controller) runWorker() {
	log.Info("Controller.runWorker: starting")

	// invoke processNextItem to fetch and consume the next change
	// to a watched or listed resource
	for c.processNextItem() {
		log.Info("Controller.runWorker: processing next item")
	}

	log.Info("Controller.runWorker: completed")
}

// processNextItem retrieves each Queued item and takes the
// necessary handler action based off of if the item was
// created or deleted
func (c *Controller) processNextItem() bool {
	log.Info("Controller.processNextItem: start")

	// fetch the next item (blocking) from the Queue to process or
	// if a shutdown is requested then return out of this to stop
	// processing
	key, quit := c.Queue.Get()

	// stop the worker loop from running as this indicates we
	// have sent a shutdown message that the Queue has indicated
	// from the Get method
	if quit {
		return false
	}

	defer c.Queue.Done(key)

	// assert the string out of the key (format `namespace/name`)
	keyRaw := key.(string)

	// take the string key and get the object out of the indexer
	//
	// item will contain the complex object for the resource and
	// exists is a bool that'll indicate whether or not the
	// resource was created (true) or deleted (false)
	//
	// if there is an error in getting the key from the index
	// then we want to retry this particular Queue key a certain
	// number of times (5 here) before we forget the Queue key
	// and throw an error
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

	// if the item doesn't exist then it was deleted and we need to fire off the handler's
	// ObjectDeleted method. but if the object does exist that indicates that the object
	// was created (or updated) so run the ObjectCreated method
	//
	// after both instances, we want to forget the key from the Queue, as this indicates
	// a code path of successful Queue key processing
	if !exists {
		c.Logger.Infof("Controller.processNextItem: object deleted detected: %s", keyRaw)
		c.Handler.ObjectDeleted(item)
		c.Queue.Forget(key)
	} else {
		c.Logger.Infof("Controller.processNextItem: object created detected: %s", keyRaw)
		c.Handler.ObjectCreated(item)
		c.Queue.Forget(key)
	}

	// keep the worker loop running by returning true
	return true
}
