package util

import (
	log "github.com/Sirupsen/logrus"
	kubebenchjobclientset "github.com/kubeflow/kubebench/controller/kubebench-operator/pkg/client/clientset/versioned"
	kubebenchjobinformer_v1 "github.com/kubeflow/kubebench/controller/kubebench-operator/pkg/client/informers/externalversions/kubebenchjob/v1"
	api_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

func GetPodsSharedIndexInformer(client kubernetes.Interface) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		// the ListWatch contains two different functions that our
		// informer requires: ListFunc to take care of listing and watching
		// the resources we want to handle
		&cache.ListWatch{
			ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
				// list all of the pods (core resource) in the deafult namespace
				return client.CoreV1().Pods(meta_v1.NamespaceDefault).List(options)
			},
			WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
				// watch all of the pods (core resource) in the default namespace
				return client.CoreV1().Pods(meta_v1.NamespaceDefault).Watch(options)
			},
		},
		&api_v1.Pod{}, // the target type (Pod)
		0,             // no resync (period of 0)
		cache.Indexers{},
	)
}

func CreateWorkingQueue() workqueue.RateLimitingInterface {
	// a result of listing or watching, we can add an idenfitying key to the queue
	// so that it can be handled in the handler
	return workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
}

func AddPodsEventHandler(inf cache.SharedInformer, queue workqueue.RateLimitingInterface) {
	// add event handlers to handle the three types of events for resources:
	//  - adding new resources
	//  - updating existing resources
	//  - deleting resources
	inf.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			// convert the resource object into a key (in this case
			// we are just doing it in the format of 'namespace/name')
			key, err := cache.MetaNamespaceKeyFunc(obj)
			log.Infof("Add pod: %s", key)
			if err == nil {
				// add the key to the queue for the handler to get
				queue.Add(key)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(newObj)
			log.Infof("Update pod: %s", key)
			if err == nil {
				queue.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			// DeletionHandlingMetaNamsespaceKeyFunc is a helper function that allows
			// us to check the DeletedFinalStateUnknown existence in the event that
			// a resource was deleted but it is still contained in the index
			//
			// this then in turn calls MetaNamespaceKeyFunc
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			log.Infof("Delete pod: %s", key)
			if err == nil {
				queue.Add(key)
			}
		},
	})
}

func GetTeamsSharedIndexInformer(client kubernetes.Interface, kubebenchjobclient kubebenchjobclientset.Interface) cache.SharedIndexInformer {
	return kubebenchjobinformer_v1.NewKubebenchJobInformer(
		kubebenchjobclient,
		meta_v1.NamespaceAll,
		0,
		cache.Indexers{},
	)
}
