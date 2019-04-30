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

package app

import (
	"os"
	"time"

	argoclientset "github.com/argoproj/argo/pkg/client/clientset/versioned"
	argoinformers "github.com/argoproj/argo/pkg/client/informers/externalversions"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	restclientset "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	kbjobv1alpha2 "github.com/kubeflow/kubebench/controller/pkg/apis/kubebenchjob/v1alpha2"
	kbjobclientset "github.com/kubeflow/kubebench/controller/pkg/client/clientset/versioned"
	kbjobinformers "github.com/kubeflow/kubebench/controller/pkg/client/informers/externalversions"
	controller "github.com/kubeflow/kubebench/controller/pkg/controller/kubebenchjob/v1alpha2"
	"github.com/kubeflow/kubebench/controller/pkg/util/signals"
)

var (
	resyncPeriod = 30 * time.Second
)

// Run creates the Kubebench controller and starts the server
func Run(opt *ServerOption) error {

	// set up signals so we handle the first shutdown signal gracefully.
	stopCh := signals.SetupSignalHandler()

	// get k8s client config
	var config *restclientset.Config
	var err error
	if opt.Kubeconfig != "" {
		config, err = clientcmd.BuildConfigFromFlags("", opt.Kubeconfig)
		if err != nil {
			log.Errorf("Failed to get config: %s", err)
			return err
		}
	} else {
		config, err = restclientset.InClusterConfig()
		if err != nil {
			log.Errorf("Failed to get config: %s", err)
			return err
		}
	}

	// get kubebench config
	kbconfig := kbjobv1alpha2.KubebenchConfig{}
	if opt.Config != "" {
		kbconfigReader, err := os.Open(opt.Config)
		if err != nil {
			log.Errorf("Failed to open kubebench config file: %s", err)
		}
		defer kbconfigReader.Close()
		if err = yaml.NewYAMLOrJSONDecoder(kbconfigReader, 128).Decode(&kbconfig); err != nil {
			log.Errorf("Failed to decode kubebench config file: %s", err)
			return err
		}
	}

	if opt.Namespace == corev1.NamespaceAll {
		log.Info("Using cluster scoped operator")
	} else {
		log.Infof("Scoping operator to namespace %s", opt.Namespace)
	}

	// create clients
	argoClientSet, kbJobClientSet, err := createClientSets(config)
	if err != nil {
		return err
	}

	// create informer factories
	argoInformerFactory := argoinformers.NewFilteredSharedInformerFactory(
		argoClientSet, resyncPeriod, opt.Namespace, nil)
	kbJobInformerFactory := kbjobinformers.NewFilteredSharedInformerFactory(
		kbJobClientSet, resyncPeriod, opt.Namespace, nil)
	// create informers
	argoInformer := argoInformerFactory.Argoproj().V1alpha1().Workflows()
	kbJobInformer := kbJobInformerFactory.Kubeflow().V1alpha2().KubebenchJobs()

	// create kubebench controller
	kbc := controller.NewKubebenchJobController(
		argoClientSet, kbJobClientSet, argoInformer, kbJobInformer, kbconfig)

	// start informer factories
	argoInformerFactory.Start(stopCh)
	kbJobInformerFactory.Start(stopCh)

	if err = kbc.Run(opt.Threadiness, stopCh); err != nil {
		log.Errorf("Failed to run Kubebench job controller: %s", err)
		return err
	}

	return nil
}

func createClientSets(config *restclientset.Config) (argoclientset.Interface, kbjobclientset.Interface, error) {
	argoClientSet, err := argoclientset.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}

	kbJobClientSet, err := kbjobclientset.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}

	return argoClientSet, kbJobClientSet, nil
}
