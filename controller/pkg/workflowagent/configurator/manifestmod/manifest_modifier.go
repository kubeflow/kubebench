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

package manifestmod

import (
	"bytes"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"

	resmod "github.com/kubeflow/kubebench/controller/pkg/resource/mod"
	"github.com/kubeflow/kubebench/controller/pkg/util"
	"github.com/kubeflow/kubebench/controller/pkg/workflowagent/configurator/common"
)

type ManifestModifier struct {
	spec *common.ManifestModSpec
}

func NewManifestModifier(modSpec *common.ManifestModSpec) *ManifestModifier {
	modifier := &ManifestModifier{spec: modSpec}
	return modifier
}

func (m *ManifestModifier) ModifyManifest(manifest []byte) ([]byte, error) {
	reader := bytes.NewReader(manifest)
	output := []byte{}
	writer := bytes.NewBuffer(output)
	resources, err := util.ReadResourcesFromYAML(reader)
	if err != nil {
		return nil, err
	}

	var modifiedResources []unstructured.Unstructured

	for _, r := range resources {
		modSpec := resmod.ResourceModSpec(*m.spec)
		modres, err := resmod.NewResourceModifier().ModifyResource(&r, &modSpec)
		if err != nil {
			return nil, err
		}
		modifiedResources = append(modifiedResources, *modres)
	}

	for i, r := range modifiedResources {
		if i > 0 {
			_, err := writer.Write([]byte("\n---\n"))
			if err != nil {
				return nil, err
			}
		}
		err = json.NewYAMLSerializer(json.DefaultMetaFactory, nil, nil).Encode(&r, writer)
		if err != nil {
			return nil, err
		}
	}

	return writer.Bytes(), nil
}
