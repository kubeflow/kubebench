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

package v1

import (
	v1 "github.com/kubeflow/kubebench/controller/kubebench-operator/pkg/apis/kubebenchjob/v1"
	scheme "github.com/kubeflow/kubebench/controller/kubebench-operator/pkg/client/clientset/versioned/scheme"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// KubebenchJobsGetter has a method to return a KubebenchJobInterface.
// A group's client should implement this interface.
type KubebenchJobsGetter interface {
	KubebenchJobs(namespace string) KubebenchJobInterface
}

// KubebenchJobInterface has methods to work with KubebenchJob resources.
type KubebenchJobInterface interface {
	Create(*v1.KubebenchJob) (*v1.KubebenchJob, error)
	Update(*v1.KubebenchJob) (*v1.KubebenchJob, error)
	UpdateStatus(*v1.KubebenchJob) (*v1.KubebenchJob, error)
	Delete(name string, options *meta_v1.DeleteOptions) error
	DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error
	Get(name string, options meta_v1.GetOptions) (*v1.KubebenchJob, error)
	List(opts meta_v1.ListOptions) (*v1.KubebenchJobList, error)
	Watch(opts meta_v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.KubebenchJob, err error)
	KubebenchJobExpansion
}

// kubebenchJobs implements KubebenchJobInterface
type kubebenchJobs struct {
	client rest.Interface
	ns     string
}

// newKubebenchJobs returns a KubebenchJobs
func newKubebenchJobs(c *KubebenchV1Client, namespace string) *kubebenchJobs {
	return &kubebenchJobs{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the kubebenchJob, and returns the corresponding kubebenchJob object, and an error if there is any.
func (c *kubebenchJobs) Get(name string, options meta_v1.GetOptions) (result *v1.KubebenchJob, err error) {
	result = &v1.KubebenchJob{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("kubebenchjobs").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of KubebenchJobs that match those selectors.
func (c *kubebenchJobs) List(opts meta_v1.ListOptions) (result *v1.KubebenchJobList, err error) {
	result = &v1.KubebenchJobList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("kubebenchjobs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested kubebenchJobs.
func (c *kubebenchJobs) Watch(opts meta_v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("kubebenchjobs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a kubebenchJob and creates it.  Returns the server's representation of the kubebenchJob, and an error, if there is any.
func (c *kubebenchJobs) Create(kubebenchJob *v1.KubebenchJob) (result *v1.KubebenchJob, err error) {
	result = &v1.KubebenchJob{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("kubebenchjobs").
		Body(kubebenchJob).
		Do().
		Into(result)
	return
}

// Update takes the representation of a kubebenchJob and updates it. Returns the server's representation of the kubebenchJob, and an error, if there is any.
func (c *kubebenchJobs) Update(kubebenchJob *v1.KubebenchJob) (result *v1.KubebenchJob, err error) {
	result = &v1.KubebenchJob{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("kubebenchjobs").
		Name(kubebenchJob.Name).
		Body(kubebenchJob).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *kubebenchJobs) UpdateStatus(kubebenchJob *v1.KubebenchJob) (result *v1.KubebenchJob, err error) {
	result = &v1.KubebenchJob{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("kubebenchjobs").
		Name(kubebenchJob.Name).
		SubResource("status").
		Body(kubebenchJob).
		Do().
		Into(result)
	return
}

// Delete takes name of the kubebenchJob and deletes it. Returns an error if one occurs.
func (c *kubebenchJobs) Delete(name string, options *meta_v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("kubebenchjobs").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *kubebenchJobs) DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("kubebenchjobs").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched kubebenchJob.
func (c *kubebenchJobs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.KubebenchJob, err error) {
	result = &v1.KubebenchJob{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("kubebenchjobs").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
