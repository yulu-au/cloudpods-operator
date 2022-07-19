/*
Copyright The Kubernetes Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"

	v1alpha1 "yunion.io/x/onecloud-operator/pkg/apis/onecloud/v1alpha1"
	scheme "yunion.io/x/onecloud-operator/pkg/client/clientset/versioned/scheme"
)

// OnecloudClustersGetter has a method to return a OnecloudClusterInterface.
// A group's client should implement this interface.
type OnecloudClustersGetter interface {
	OnecloudClusters(namespace string) OnecloudClusterInterface
}

// OnecloudClusterInterface has methods to work with OnecloudCluster resources.
type OnecloudClusterInterface interface {
	Create(*v1alpha1.OnecloudCluster) (*v1alpha1.OnecloudCluster, error)
	Update(*v1alpha1.OnecloudCluster) (*v1alpha1.OnecloudCluster, error)
	UpdateStatus(*v1alpha1.OnecloudCluster) (*v1alpha1.OnecloudCluster, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.OnecloudCluster, error)
	List(opts v1.ListOptions) (*v1alpha1.OnecloudClusterList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.OnecloudCluster, err error)
	OnecloudClusterExpansion
}

// onecloudClusters implements OnecloudClusterInterface
type onecloudClusters struct {
	client rest.Interface
	ns     string
}

// newOnecloudClusters returns a OnecloudClusters
func newOnecloudClusters(c *OnecloudV1alpha1Client, namespace string) *onecloudClusters {
	return &onecloudClusters{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the onecloudCluster, and returns the corresponding onecloudCluster object, and an error if there is any.
func (c *onecloudClusters) Get(name string, options v1.GetOptions) (result *v1alpha1.OnecloudCluster, err error) {
	result = &v1alpha1.OnecloudCluster{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("onecloudclusters").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of OnecloudClusters that match those selectors.
func (c *onecloudClusters) List(opts v1.ListOptions) (result *v1alpha1.OnecloudClusterList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.OnecloudClusterList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("onecloudclusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested onecloudClusters.
func (c *onecloudClusters) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("onecloudclusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a onecloudCluster and creates it.  Returns the server's representation of the onecloudCluster, and an error, if there is any.
func (c *onecloudClusters) Create(onecloudCluster *v1alpha1.OnecloudCluster) (result *v1alpha1.OnecloudCluster, err error) {
	result = &v1alpha1.OnecloudCluster{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("onecloudclusters").
		Body(onecloudCluster).
		Do().
		Into(result)
	return
}

// Update takes the representation of a onecloudCluster and updates it. Returns the server's representation of the onecloudCluster, and an error, if there is any.
func (c *onecloudClusters) Update(onecloudCluster *v1alpha1.OnecloudCluster) (result *v1alpha1.OnecloudCluster, err error) {
	result = &v1alpha1.OnecloudCluster{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("onecloudclusters").
		Name(onecloudCluster.Name).
		Body(onecloudCluster).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *onecloudClusters) UpdateStatus(onecloudCluster *v1alpha1.OnecloudCluster) (result *v1alpha1.OnecloudCluster, err error) {
	result = &v1alpha1.OnecloudCluster{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("onecloudclusters").
		Name(onecloudCluster.Name).
		SubResource("status").
		Body(onecloudCluster).
		Do().
		Into(result)
	return
}

// Delete takes name of the onecloudCluster and deletes it. Returns an error if one occurs.
func (c *onecloudClusters) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("onecloudclusters").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *onecloudClusters) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("onecloudclusters").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched onecloudCluster.
func (c *onecloudClusters) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.OnecloudCluster, err error) {
	result = &v1alpha1.OnecloudCluster{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("onecloudclusters").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
