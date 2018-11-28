/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package fake

import (
	kubebenchjob_v1 "github.com/kubeflow/kubebench/controller/kubebench-operator/pkg/apis/kubebenchjob/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeKubebenchJobs implements KubebenchJobInterface
type FakeKubebenchJobs struct {
	Fake *FakeKubebenchV1
	ns   string
}

var kubebenchjobsResource = schema.GroupVersionResource{Group: "kubebench.operator", Version: "v1", Resource: "kubebenchjobs"}

var kubebenchjobsKind = schema.GroupVersionKind{Group: "kubebench.operator", Version: "v1", Kind: "KubebenchJob"}

// Get takes name of the kubebenchJob, and returns the corresponding kubebenchJob object, and an error if there is any.
func (c *FakeKubebenchJobs) Get(name string, options v1.GetOptions) (result *kubebenchjob_v1.KubebenchJob, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(kubebenchjobsResource, c.ns, name), &kubebenchjob_v1.KubebenchJob{})

	if obj == nil {
		return nil, err
	}
	return obj.(*kubebenchjob_v1.KubebenchJob), err
}

// List takes label and field selectors, and returns the list of KubebenchJobs that match those selectors.
func (c *FakeKubebenchJobs) List(opts v1.ListOptions) (result *kubebenchjob_v1.KubebenchJobList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(kubebenchjobsResource, kubebenchjobsKind, c.ns, opts), &kubebenchjob_v1.KubebenchJobList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &kubebenchjob_v1.KubebenchJobList{}
	for _, item := range obj.(*kubebenchjob_v1.KubebenchJobList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested kubebenchJobs.
func (c *FakeKubebenchJobs) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(kubebenchjobsResource, c.ns, opts))

}

// Create takes the representation of a kubebenchJob and creates it.  Returns the server's representation of the kubebenchJob, and an error, if there is any.
func (c *FakeKubebenchJobs) Create(kubebenchJob *kubebenchjob_v1.KubebenchJob) (result *kubebenchjob_v1.KubebenchJob, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(kubebenchjobsResource, c.ns, kubebenchJob), &kubebenchjob_v1.KubebenchJob{})

	if obj == nil {
		return nil, err
	}
	return obj.(*kubebenchjob_v1.KubebenchJob), err
}

// Update takes the representation of a kubebenchJob and updates it. Returns the server's representation of the kubebenchJob, and an error, if there is any.
func (c *FakeKubebenchJobs) Update(kubebenchJob *kubebenchjob_v1.KubebenchJob) (result *kubebenchjob_v1.KubebenchJob, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(kubebenchjobsResource, c.ns, kubebenchJob), &kubebenchjob_v1.KubebenchJob{})

	if obj == nil {
		return nil, err
	}
	return obj.(*kubebenchjob_v1.KubebenchJob), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeKubebenchJobs) UpdateStatus(kubebenchJob *kubebenchjob_v1.KubebenchJob) (*kubebenchjob_v1.KubebenchJob, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(kubebenchjobsResource, "status", c.ns, kubebenchJob), &kubebenchjob_v1.KubebenchJob{})

	if obj == nil {
		return nil, err
	}
	return obj.(*kubebenchjob_v1.KubebenchJob), err
}

// Delete takes name of the kubebenchJob and deletes it. Returns an error if one occurs.
func (c *FakeKubebenchJobs) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(kubebenchjobsResource, c.ns, name), &kubebenchjob_v1.KubebenchJob{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeKubebenchJobs) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(kubebenchjobsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &kubebenchjob_v1.KubebenchJobList{})
	return err
}

// Patch applies the patch and returns the patched kubebenchJob.
func (c *FakeKubebenchJobs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *kubebenchjob_v1.KubebenchJob, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(kubebenchjobsResource, c.ns, name, data, subresources...), &kubebenchjob_v1.KubebenchJob{})

	if obj == nil {
		return nil, err
	}
	return obj.(*kubebenchjob_v1.KubebenchJob), err
}
