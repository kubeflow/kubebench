package main

import (
	log "github.com/Sirupsen/logrus"

	"os"
	"os/signal"
	"syscall"

	"github.com/kubeflow/kubebench/controller/kubebench-operator/pkg/client"
	controllers "github.com/kubeflow/kubebench/controller/kubebench-operator/pkg/controller"
	"github.com/kubeflow/kubebench/controller/kubebench-operator/pkg/handler"
	"github.com/kubeflow/kubebench/controller/kubebench-operator/pkg/util"
)

func main() {
	// Get the Kubernetes client to access the Cloud platform
	client, kubebenchjobclient := client.GetKubernetesCRDClient()

	teaminformer := util.GetTeamsSharedIndexInformer(client, kubebenchjobclient)
	queue := util.CreateWorkingQueue()
	util.AddPodsEventHandler(teaminformer, queue)

	// construct the Controller object which has all of the necessary components to
	// handle logging, connections, informing (listing and watching), the queue,
	// and the handler
	controller := controllers.KubebenchJobController{
		Logger:    log.NewEntry(log.New()),
		Clientset: client,
		Informer:  teaminformer,
		Queue:     queue,
		Handler:   handler.KubebenchJobHandler{},
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
