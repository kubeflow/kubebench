package client

import (
	"flag"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	kubebenchjobclientset "github.com/kubeflow/kubebench/controller/kubebench-operator/pkg/client/clientset/versioned"
)

var config = parseKubernetesConfig()

func parseKubernetesConfig() *restclient.Config {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	// Parse the command line arguments
	flag.Parse()

	// create the config from the path
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Fatalf("getClusterConfig: %v", err)
	}
	return config
}

// Retrieve the Kubernetes cluster client from outside of the cluster and add the Team Clienset
func GetKubernetesCRDClient() (kubernetes.Interface, kubebenchjobclientset.Interface) {
	// Generate the client based off of the config
	client := GetKubernetesClient()

	// Create a Team ClientSet
	clientset, err := kubebenchjobclientset.NewForConfig(config)
	if err != nil {
		log.Fatalf("Team clienset: %v", err)
	}

	log.Info("Successfully constructed k8s client")
	return client, clientset
}

// Retrieve the Kubernetes cluster client from outside of the cluster
func GetKubernetesClient() kubernetes.Interface {

	// generate the client based off of the config
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("getClusterConfig: %v", err)
	}

	log.Info("Successfully constructed k8s client")
	return client
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
