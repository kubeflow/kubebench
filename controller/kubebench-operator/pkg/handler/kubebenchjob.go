package handler

import (
	log "github.com/Sirupsen/logrus"
	kubebenchjob_v1 "github.com/kubeflow/kubebench/controller/kubebench-operator/pkg/apis/kubebenchjob/v1"
)

// Handler interface contains the methods that are required
type HandlerInterface interface {
	Init() error
	ObjectCreated(obj interface{})
	ObjectDeleted(obj interface{})
	ObjectUpdated(objOld, objNew interface{})
}

type KubebenchJobHandler struct{}

func (t *KubebenchJobHandler) Init() error {
	log.Info("TeamHandler.Init")
	return nil
}

// ObjectCreated is called when an object is created
func (t *KubebenchJobHandler) ObjectCreated(obj interface{}) {
	log.Info("TeamHandler.ObjectCreated")

	kbjob := obj.(*kubebenchjob_v1.KubebenchJob)
	log.Infof("    ResourceVersion: %s", kbjob.ObjectMeta.Name)
}

// ObjectDeleted is called when an object is deleted
func (t *KubebenchJobHandler) ObjectDeleted(obj interface{}) {
	log.Info("KubebenchJobHandler.ObjectDeleted")
}

// ObjectUpdated is called when an object is updated
func (t *KubebenchJobHandler) ObjectUpdated(objOld, objNew interface{}) {
	log.Info("KubebenchJobHandler.ObjectUpdated")
}
