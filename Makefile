verify:
	go fmt ./controller/cmd/configurator/
	go fmt ./controller/cmd/reporter/
	go fmt ./controller/pkg/apis/kubebench/v1alpha1/
	go fmt ./controller/pkg/util/