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
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"time"

	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/rest"

	"github.com/kubeflow/kubebench/controller/pkg/resource/client"
	"github.com/kubeflow/kubebench/controller/pkg/resource/common"
)

// ResourceManager applies high level actions to resources
type ResourceManager struct {
	clusterConfig *rest.Config
	inputReader   io.Reader
	outputWriter  io.Writer
}

// NewResourceManager creates a new ResourceManager for the given config
func NewResourceManager(config *rest.Config, reader io.Reader, writer io.Writer) *ResourceManager {
	resourceManager := ResourceManager{
		clusterConfig: config,
		inputReader:   reader,
		outputWriter:  writer,
	}
	return &resourceManager
}

// Run executes ResourceManager with given options
func (m *ResourceManager) Run(opt *AppOption) error {
	resourceClient, err := client.NewResourceClient(m.clusterConfig)
	if err != nil {
		log.Errorf("Failed to create resource client: %s", err)
		return err
	}
	switch opt.Action {
	case "create":
		resourceObjects, err := getResourceObjects(m.inputReader)
		if err != nil {
			log.Errorf("Failed to get resource objects: %s", err)
			return err
		}
		results, err := resourceClient.Create(resourceObjects, opt.NumCopies)
		if err != nil {
			log.Errorf("Failed to run resource client: %s", err)
			return err
		}
		if err := writeCreateResults(m.outputWriter, results); err != nil {
			log.Errorf("Failed to write create results: %s", err)
			return err
		}
	case "auto-watch":
		resourceRefs, err := getResourceRefs(m.inputReader)
		if err != nil {
			log.Errorf("Failed to get resource-refs: %s", err)
			return err
		}
		timeout, err := time.ParseDuration(opt.Timeout)
		if err != nil {
			log.Errorf("Failed to parse timeout value: %s", err)
		}
		results, err := resourceClient.AutoWatchByRef(resourceRefs, timeout)
		if err != nil {
			log.Errorf("Failed to run resource client: %s", err)
			return err
		}
		if err := writeAutoWatchResults(m.outputWriter, results); err != nil {
			log.Errorf("Failed to write auto-watch results: %s", err)
			return err
		}
	default:
		err := fmt.Errorf("Unknown action: %s", opt.Action)
		log.Errorf("s", err)
		return err
	}
	return nil
}

func getResourceObjects(reader io.Reader) ([]*unstructured.Unstructured, error) {
	yamlReader := yaml.NewYAMLReader(bufio.NewReader(reader))
	var resources []*unstructured.Unstructured

	for {
		data, err := yamlReader.Read()
		if err != nil && err != io.EOF {
			return nil, err
		} else if err == io.EOF {
			break
		}

		data, err = yaml.ToJSON(data)
		if err != nil {
			return nil, err
		}

		obj := unstructured.Unstructured{}
		if err := json.Unmarshal(data, &obj); err != nil {
			log.Errorf("Failed to parse manifest to k8s object: %s", err)
			return nil, err
		}
		resources = append(resources, &obj)
	}

	return resources, nil
}

func getResourceRefs(reader io.Reader) ([]*common.ResourceRef, error) {
	jsonDecoder := json.NewDecoder(reader)
	var resourceRefs []*common.ResourceRef

	for {
		var ref common.ResourceRef
		if err := jsonDecoder.Decode(&ref); err == io.EOF {
			break
		} else if err != nil {
			log.Errorf("Failed to parse resource-ref file: %s", err)
			return nil, err
		}
		resourceRefs = append(resourceRefs, &ref)
	}

	return resourceRefs, nil
}

func writeCreateResults(writer io.Writer, results []*common.ResourceRef) error {
	jsonEncoder := json.NewEncoder(writer)
	for _, ref := range results {
		if err := jsonEncoder.Encode(ref); err != nil {
			log.Errorf("Failed to write resource-ref file: %s", err)
			return err
		}
	}
	return nil
}

func writeAutoWatchResults(writer io.Writer, results []*client.WatchResult) error {
	jsonEncoder := json.NewEncoder(writer)
	for _, ref := range results {
		if err := jsonEncoder.Encode(ref); err != nil {
			log.Errorf("Failed to write watch-result file: %s", err)
			return err
		}
	}
	return nil
}
