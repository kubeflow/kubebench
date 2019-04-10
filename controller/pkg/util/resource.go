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

package util

import (
	"bufio"
	"encoding/json"
	"io"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// ReadResourcesFromYAML reads in k8s resources as a list of Unstructured from YAML input
func ReadResourcesFromYAML(reader io.Reader) ([]unstructured.Unstructured, error) {
	yamlReader := yaml.NewYAMLReader(bufio.NewReader(reader))
	var resources []unstructured.Unstructured

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
			return nil, err
		}
		resources = append(resources, obj)
	}

	return resources, nil
}
