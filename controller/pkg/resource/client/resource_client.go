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

package client

import (
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"

	"github.com/kubeflow/kubebench/controller/pkg/resource/common"
	"github.com/kubeflow/kubebench/controller/pkg/resource/condition"
)

// ResourceClient interacts with the cluster and performs resource operations
type ResourceClient struct {
	discoveryClient *discovery.DiscoveryClient
	dynamicClient   dynamic.Interface
}

// NewResourceClient creates a new ResourceClient
func NewResourceClient(config *rest.Config) (*ResourceClient, error) {
	discc, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		log.Errorf("Failed to create discovery client: %s", err)
		return nil, err
	}
	dynac, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Errorf("Failed to create dynamic client: %s", err)
		return nil, err
	}
	return &ResourceClient{discoveryClient: discc, dynamicClient: dynac}, nil
}

// Create creates a list of resources to k8s cluster
func (rc *ResourceClient) Create(resources []*unstructured.Unstructured, numCopies int) ([]*common.ResourceRef, error) {
	var results []*common.ResourceRef
	for i := 0; i < numCopies; i++ {
		for _, r := range resources {
			rcopy := r.DeepCopy()
			if name := rcopy.GetName(); numCopies > 1 && name != "" {
				rcopy.SetGenerateName(name + "-")
				rcopy.SetName("")
			}
			ref := common.NewResourceRefFromUnstructured(rcopy)
			kvgn := ref.SprintKindVersionGroupName()
			ns := ref.Namespace
			log.Infof("Creating resource: %s in namespace: %s", kvgn, ns)
			apiRes, err := rc.getAPIResource(ref)
			if err != nil {
				return results, err
			}
			gvr := schema.GroupVersionResource{
				Group:    ref.Group,
				Version:  ref.Version,
				Resource: apiRes.Name,
			}
			var resIf dynamic.ResourceInterface
			if apiRes.Namespaced {
				resIf = rc.dynamicClient.Resource(gvr).Namespace(ns)
			} else {
				resIf = rc.dynamicClient.Resource(gvr)
			}
			res, err := resIf.Create(rcopy, metav1.CreateOptions{})
			if err != nil {
				log.Errorf("Failed to create resource %s: %s", kvgn, err)
				return results, err
			}
			results = append(results, common.NewResourceRefFromUnstructured(res))
		}
	}
	return results, nil
}

// AutoWatchByRef automatically watches a list of resources identified by ResourceRef.
// Only supported types of resources are watched, others are ignored.
func (rc *ResourceClient) AutoWatchByRef(refs []*common.ResourceRef, timeout time.Duration) ([]*WatchResult, error) {
	var resultList []*WatchResult
	var watchList []watchItem

	for _, ref := range refs {
		resCond := condition.NewResourceCondition(ref)
		if resCond != nil {
			apiRes, err := rc.getAPIResource(ref)
			if err != nil {
				return nil, err
			}
			item := watchItem{resourceRef: ref, apiResource: apiRes, resourceCond: resCond}
			watchList = append(watchList, item)
		}
	}

	resultChan := make(chan *WatchResult, len(watchList))
	for _, wi := range watchList {
		go func(resultChan chan<- *WatchResult) {
			backoff := condition.Backoff{
				Duration: time.Second * 5,
				Factor:   2.0,
				Jitter:   0.2,
				Cap:      time.Minute * 10,
			}
			var resCondStatus condition.ResourceConditionStatus
			err := condition.ExponentialBackoff(backoff, timeout, func() (bool, error) {
				resObj, err := rc.getOneByRef(wi.resourceRef, wi.apiResource)
				if err != nil {
					return true, err
				}
				resCondStatus, err = wi.resourceCond.CheckCondition(resObj)
				if err != nil {
					return true, err
				}
				if resCondStatus == condition.ResourceConditionSuccess ||
					resCondStatus == condition.ResourceConditionFailure {
					return true, nil
				}
				return false, nil
			})
			if err != nil {
				if err == condition.ErrWaitTimeout {
					resCondStatus = condition.ResourceConditionTimeout
					log.Errorf("Waiting for resource %s resulted in timeout",
						wi.resourceRef.SprintKindVersionGroupName())
				} else {
					log.Errorf("Waiting for resource %s resulted in error: %s",
						wi.resourceRef.SprintKindVersionGroupName(), err)
				}
			}
			resultChan <- &WatchResult{ResourceRef: wi.resourceRef, Status: resCondStatus, Error: err}
		}(resultChan)
	}

	for i := 0; i < len(watchList); i++ {
		result := <-resultChan
		resultList = append(resultList, result)
	}

	return resultList, nil
}

func (rc *ResourceClient) getAPIResource(r *common.ResourceRef) (*metav1.APIResource, error) {
	apiResourceLists, err := rc.discoveryClient.ServerResources()
	if err != nil {
		return nil, err
	}
	var result *metav1.APIResource
	// find the resource name from "kind"
	for _, resList := range apiResourceLists {
		groupVersion := resList.GroupVersion
		if groupVersion != r.SprintGroupVersion() {
			continue
		}
		for _, res := range resList.APIResources {
			// find the resource that matches "kind", and is not a subresource
			if res.Kind == r.Kind && !strings.Contains(res.Name, "/") {
				result = &res
				break
			}
		}
	}
	if result == nil {
		err = fmt.Errorf("Cannot find resource for %s", r.SprintKindVersionGroup())
		return nil, err
	}

	return result, nil
}

// getOneByRef gets single resource's info from k8s cluster
func (rc *ResourceClient) getOneByRef(ref *common.ResourceRef, apiRes *metav1.APIResource) (*unstructured.Unstructured, error) {
	gvr := schema.GroupVersionResource{
		Group:    ref.Group,
		Version:  ref.Version,
		Resource: apiRes.Name,
	}
	var resIf dynamic.ResourceInterface
	if apiRes.Namespaced {
		resIf = rc.dynamicClient.Resource(gvr).Namespace(ref.Namespace)
	} else {
		resIf = rc.dynamicClient.Resource(gvr)
	}
	res, err := resIf.Get(ref.Name, metav1.GetOptions{})
	if err != nil {
		log.Errorf("Failed to get resource %s: %s", ref.SprintKindVersionGroupName(), err)
		return nil, err
	}
	return res, nil
}

type watchItem struct {
	resourceRef  *common.ResourceRef
	apiResource  *metav1.APIResource
	resourceCond condition.ResourceConditionInterface
}

// WatchResult is the result of auto-watch
type WatchResult struct {
	ResourceRef *common.ResourceRef               `json:"resourceRef"`
	Status      condition.ResourceConditionStatus `json:"status"`
	Error       error                             `json:"-"`
}
