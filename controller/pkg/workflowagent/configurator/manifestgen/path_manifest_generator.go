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

package manifestgen

import (
	"io/ioutil"
)

type PathManifestGenerator struct {
	spec *string
}

func NewPathManifestGenerator(spec interface{}) ManifestGeneratorInterface {
	pathSpec := spec.(*string)
	generator := &PathManifestGenerator{
		spec: pathSpec,
	}
	return generator
}

func (mg *PathManifestGenerator) GenerateManifest() ([]byte, error) {
	path := mg.spec
	data, err := ioutil.ReadFile(string(*path))
	if err != nil {
		return nil, err
	}
	return data, nil
}
