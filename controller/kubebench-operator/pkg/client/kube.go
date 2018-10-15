package client

import (
	"flag"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	argoproj "github.com/argoproj/argo/pkg/client/clientset/versioned/typed/workflow/v1alpha1"
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
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Fatalf("getClusterConfig: %v", err)
	}
	return config
}

func GetKubernetesCRDClient() (kubernetes.Interface, kubebenchjobclientset.Interface) {
	client := GetKubernetesClient()

	clientset, err := kubebenchjobclientset.NewForConfig(config)
	if err != nil {
		log.Fatalf("KubebenchJob clienset: %v", err)
	}

	log.Info("Successfully constructed k8s client")
	return client, clientset
}

func GetArgoClient() *argoproj.ArgoprojV1alpha1Client {
	argoClient, err := argoproj.NewForConfig(config)
	if err != nil {
		log.Fatalf("ArgoClient: %v", err)
	}

	return argoClient
}

func GetKubernetesClient() kubernetes.Interface {

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
