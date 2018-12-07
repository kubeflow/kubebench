package client

import (
	"flag"
	"os"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	argoproj "github.com/argoproj/argo/pkg/client/clientset/versioned/typed/workflow/v1alpha1"
	kubebenchjobclientset "github.com/kubeflow/kubebench/controller/pkg/client/clientset/versioned"
)

var config = parseKubernetesConfig()

func parseKubernetesConfig() *restclient.Config {
	// if home := homeDir(); home != "" {
	// 	kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	// } else {
	// 	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	// }
	kubeconfig := flag.String("kubeconfig", "", "Path to a kube config. Only required if out-of-cluster.")
	flag.Parse()

	var err error
	var config *restclient.Config
	if *kubeconfig != "" {
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
	} else {
		config, err = restclient.InClusterConfig()
	}

	// config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Fatalf("getClusterConfig: %v", err)
	}

	// namespace, err := config.Namespace()
	// if err != nil {
	// 	log.Fatalf("Error getting namespace: %v", err)
	// }

	// log.Infof("Namespace is %s", namespace)

	// fooType := reflect.TypeOf(config)
	// for i := 0; i < fooType.NumMethod(); i++ {
	// 	method := fooType.Method(i)
	// 	fmt.Println(method.Name)
	// }
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
