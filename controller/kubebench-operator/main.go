package main

// func main() {
// 	// Get the Kubernetes client to access the Cloud platform
// 	client := client.GetKubernetesClient()

// 	ns, nsError := client.CoreV1().Namespaces().List(metav1.ListOptions{})
// 	if nsError != nil {
// 		log.Fatalf("Can't list namespaces ", nsError)
// 	}
// 	for i := range ns.Items {
// 		log.Info("Namespace/project : ", ns.Items[i].Name)
// 	}
// }
