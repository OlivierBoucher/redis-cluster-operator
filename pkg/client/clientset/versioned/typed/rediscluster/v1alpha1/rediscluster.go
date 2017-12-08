package v1alpha1

import (
	v1alpha1 "github.com/OlivierBoucher/redis-cluster-operator/pkg/apis/rediscluster/v1alpha1"
	scheme "github.com/OlivierBoucher/redis-cluster-operator/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// RedisClustersGetter has a method to return a RedisClusterInterface.
// A group's client should implement this interface.
type RedisClustersGetter interface {
	RedisClusters(namespace string) RedisClusterInterface
}

// RedisClusterInterface has methods to work with RedisCluster resources.
type RedisClusterInterface interface {
	Create(*v1alpha1.RedisCluster) (*v1alpha1.RedisCluster, error)
	Update(*v1alpha1.RedisCluster) (*v1alpha1.RedisCluster, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.RedisCluster, error)
	List(opts v1.ListOptions) (*v1alpha1.RedisClusterList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.RedisCluster, err error)
	RedisClusterExpansion
}

// redisClusters implements RedisClusterInterface
type redisClusters struct {
	client rest.Interface
	ns     string
}

// newRedisClusters returns a RedisClusters
func newRedisClusters(c *RedisclusterV1alpha1Client, namespace string) *redisClusters {
	return &redisClusters{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the redisCluster, and returns the corresponding redisCluster object, and an error if there is any.
func (c *redisClusters) Get(name string, options v1.GetOptions) (result *v1alpha1.RedisCluster, err error) {
	result = &v1alpha1.RedisCluster{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("redisclusters").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of RedisClusters that match those selectors.
func (c *redisClusters) List(opts v1.ListOptions) (result *v1alpha1.RedisClusterList, err error) {
	result = &v1alpha1.RedisClusterList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("redisclusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested redisClusters.
func (c *redisClusters) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("redisclusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a redisCluster and creates it.  Returns the server's representation of the redisCluster, and an error, if there is any.
func (c *redisClusters) Create(redisCluster *v1alpha1.RedisCluster) (result *v1alpha1.RedisCluster, err error) {
	result = &v1alpha1.RedisCluster{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("redisclusters").
		Body(redisCluster).
		Do().
		Into(result)
	return
}

// Update takes the representation of a redisCluster and updates it. Returns the server's representation of the redisCluster, and an error, if there is any.
func (c *redisClusters) Update(redisCluster *v1alpha1.RedisCluster) (result *v1alpha1.RedisCluster, err error) {
	result = &v1alpha1.RedisCluster{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("redisclusters").
		Name(redisCluster.Name).
		Body(redisCluster).
		Do().
		Into(result)
	return
}

// Delete takes name of the redisCluster and deletes it. Returns an error if one occurs.
func (c *redisClusters) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("redisclusters").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *redisClusters) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("redisclusters").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched redisCluster.
func (c *redisClusters) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.RedisCluster, err error) {
	result = &v1alpha1.RedisCluster{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("redisclusters").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
