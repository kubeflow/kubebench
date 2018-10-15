package main

import (
	log "github.com/Sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kubeflow/kubebench/controller/kubebench-operator/pkg/client"
	controllers "github.com/kubeflow/kubebench/controller/kubebench-operator/pkg/controller"
	"github.com/kubeflow/kubebench/controller/kubebench-operator/pkg/handler"
	"github.com/kubeflow/kubebench/controller/kubebench-operator/pkg/util"

	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Get the Kubernetes client to access the Cloud platform
	client := client.GetKubernetesClient()

	ns, nsError := client.CoreV1().Namespaces().List(metav1.ListOptions{})
	if nsError != nil {
		log.Fatalf("Can't list namespaces ", nsError)
	}
	for i := range ns.Items {
		log.Info("Namespace/project : ", ns.Items[i].Name)
	}
	informer := util.GetPodsSharedIndexInformer(client)
	queue := util.CreateWorkingQueue()
	util.AddPodsEventHandler(informer, queue)

	// construct the Controller object which has all of the necessary components to
	// handle logging, connections, informing (listing and watching), the queue,
	// and the handler
	controller := controllers.Controller{
		Logger:    log.NewEntry(log.New()),
		Clientset: client,
		Informer:  informer,
		Queue:     queue,
		Handler:   handler.SimpleHandler{},
	}

	// use a channel to synchronize the finalization for a graceful shutdown
	stopCh := make(chan struct{})
	defer close(stopCh)

	// run the controller loop to process items
	go controller.Run(stopCh)

	// use a channel to handle OS signals to terminate and gracefully shut
	// down processing
	sigTerm := make(chan os.Signal, 1)
	signal.Notify(sigTerm, syscall.SIGTERM)
	signal.Notify(sigTerm, syscall.SIGINT)
	<-sigTerm
}
