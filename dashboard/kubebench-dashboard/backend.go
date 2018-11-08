package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	argoproj "github.com/argoproj/argo/pkg/client/clientset/versioned/typed/workflow/v1alpha1"
	"github.com/ghodss/yaml"
	kbjob "github.com/kubeflow/kubebench/controller/kubebench-operator/pkg/apis/kubebenchjob/v1"
	kubeclient "github.com/kubeflow/kubebench/controller/kubebench-operator/pkg/client"
	kubebenchjobclientset "github.com/kubeflow/kubebench/controller/kubebench-operator/pkg/client/clientset/versioned"
	kubebench "github.com/kubeflow/kubebench/controller/kubebench-operator/pkg/client/clientset/versioned/typed/kubebenchjob/v1"
	"github.com/kubeflow/kubebench/controller/kubebench-operator/pkg/util"
	utils "github.com/kubeflow/kubebench/controller/kubebench-operator/pkg/util"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
)

var (
	port           = "9303"
	allowedHeaders = "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token"
)

var config = parseKubernetesConfig()

var (
	workflows     argoproj.WorkflowInterface
	kubebenchJobs kubebench.KubebenchJobInterface
)

func parseKubernetesConfig() *restclient.Config {

	config, err := restclient.InClusterConfig()
	if err != nil {
		log.Fatalf("getClusterConfig: %v", err)
	}
	return config
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

func GetKubernetesCRDClient() (kubernetes.Interface, kubebenchjobclientset.Interface) {
	client := GetKubernetesClient()

	clientset, err := kubebenchjobclientset.NewForConfig(config)
	if err != nil {
		log.Fatalf("KubebenchJob clienset: %v", err)
	}

	log.Info("Successfully constructed k8s client")
	return client, clientset
}

func main() {

	client, kubebenchjobclient := kubeclient.GetKubernetesCRDClient()

	teaminformer := util.GetTeamsSharedIndexInformer(client, kubebenchjobclient)
	queue := util.CreateWorkingQueue()
	util.AddPodsEventHandler(teaminformer, queue)

	// get kb job lclient

	// get newKubebenchJobs
	// namespace

	kbJobClient, err := kubebench.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to start kubebench client: %v", err)
	}

	argoClient, err := argoproj.NewForConfig(config)
	if err != nil {
		log.Fatalf("ArgoClient: %v", err)
	}

	workflows = argoClient.Workflows("default")

	kubebenchJobs = kbJobClient.KubebenchJobs("default")

	frontend := http.FileServer(http.Dir("/go/src/github.com/kubeflow/kubebench/dashboard/kubebench-dashboard/frontend/build/"))
	http.Handle("/dashboard/", http.StripPrefix("/dashboard/", frontend))
	// http.HandleFunc("/", mainHandler)
	http.HandleFunc("/dashboard/submit_yaml/", submitYamlHandler)
	http.HandleFunc("/dashboard/submit_params/", submitParamHandler)
	http.HandleFunc("/dashboard/fetch_jobs/", fetchJobsHandler)
	http.HandleFunc("/dashboard/delete_job/", deleteJobHandler)
	log.Println("Listening on", ":"+port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", allowedHeaders)
	(*w).Header().Set("Access-Control-Expose-Headers", "Authorization")
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func submitYamlHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	var data map[string]interface{}

	json.NewDecoder(r.Body).Decode(&data)

	job := kbjob.KubebenchJob{}
	if yamlContent, ok := data["yaml"].(string); ok {
		err := yaml.Unmarshal([]byte(yamlContent), &job)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = kubebenchJobs.Create(&job)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func submitParamHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	var data map[string]interface{}

	json.NewDecoder(r.Body).Decode(&data)

	parameters := make(map[string]string)
	for key, value := range data {
		strKey := fmt.Sprintf("%v", key)
		strValue := fmt.Sprintf("%v", value)
		parameters[strKey] = strValue
	}
	if len(parameters) != 0 {
		job, err := utils.GenerateJobFromParameters(parameters)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//		fmt.Println(job)
		_, err = kubebenchJobs.Create(job)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

type ReturnJobs struct {
	Names  []string
	Status []string
}

func fetchJobsHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	returnJobs := ReturnJobs{}

	options := metav1.ListOptions{}
	options.TypeMeta = metav1.TypeMeta{
		Kind: "KubebenchJob",
	}
	jobs, err := kubebenchJobs.List(options)
	if err != nil {
		log.Fatalf("Failed to list job: %v", err)
	}

	names := make([]string, 0)
	status := make([]string, 0)

	for _, job := range jobs.Items {
		names = append(names, job.ObjectMeta.Name)
		workflow, err := workflows.Get(job.ObjectMeta.Name, metav1.GetOptions{})
		if err != nil {
			log.Fatalf("Failed to list job: %v", err)
		}
		status = append(status, string(workflow.Status.Phase))

	}
	returnJobs.Names = names
	returnJobs.Status = status
	response, err := json.Marshal(returnJobs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

type DeleteJob struct {
	Status string
}

func deleteJobHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	deleteJob := DeleteJob{}
	var data map[string]interface{}

	json.NewDecoder(r.Body).Decode(&data)
	if data["name"] != nil {
		if name, ok := data["name"].(string); ok {
			if name != "" {
				// delete from indexer
				err := kubebenchJobs.Delete(name, &metav1.DeleteOptions{})
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
			/* act on str */
		}
	}

	deleteJob.Status = "ok"
	response, err := json.Marshal(deleteJob)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}
