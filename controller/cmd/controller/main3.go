package main

import (
	"os"
	"os/signal"
	"syscall"

	kubeclient "github.com/kubeflow/kubebench/controller/pkg/client"
	controllers "github.com/kubeflow/kubebench/controller/pkg/controller"
	log "github.com/sirupsen/logrus"

	// "github.com/kubeflow/kubebench/controller/pkg/handler"
	"github.com/kubeflow/kubebench/controller/pkg/util"
)

func main() {
	client, kubebenchjobclient := kubeclient.GetKubernetesCRDClient()

	teaminformer := util.GetTeamsSharedIndexInformer(client, kubebenchjobclient)
	queue := util.CreateWorkingQueue()
	util.AddPodsEventHandler(teaminformer, queue)

	argoClient := kubeclient.GetArgoClient()

	controller := controllers.KubebenchJobController{
		Logger:    log.NewEntry(log.New()),
		Clientset: client,
		Informer:  teaminformer,
		Queue:     queue,
		// Handler:   handler.KubebenchJobHandler{},

		//pass correct namespace here
		Workflows: argoClient.Workflows("default"),
	}

	stopCh := make(chan struct{})
	defer close(stopCh)

	go controller.Run(stopCh)

	// use a channel to handle OS signals to terminate and gracefully shut
	// down processing
	sigTerm := make(chan os.Signal, 1)
	signal.Notify(sigTerm, syscall.SIGTERM)
	signal.Notify(sigTerm, syscall.SIGINT)
	<-sigTerm
}
