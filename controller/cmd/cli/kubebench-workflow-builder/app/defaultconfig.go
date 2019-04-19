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

// default kubebench config
// NOTE: indent using spaces to form a valid yaml doc
var defaultConfig = `
defaultWorkflowAgent:
  container:
    name: kubebench-workflow-agent
    image: ciscoai/kubebench-workflow-agent:latest
defaultManagedVolumes:
  experimentVolume:
    name: kubebench-experiment-volume
    emptyDir: {}
  workflowVolume:
    name: kubebench-workflow-volume
    emptyDir: {}
`
