package fake

import (
	v1alpha1 "github.com/OlivierBoucher/redis-cluster-operator/pkg/apis/rediscluster/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeRedisClusters implements RedisClusterInterface
type FakeRedisClusters struct {
	Fake *FakeRedisclusterV1alpha1
	ns   string
}

var redisclustersResource = schema.GroupVersionResource{Group: "rediscluster.everflow.io", Version: "v1alpha1", Resource: "redisclusters"}

var redisclustersKind = schema.GroupVersionKind{Group: "rediscluster.everflow.io", Version: "v1alpha1", Kind: "RedisCluster"}

// Get takes name of the redisCluster, and returns the corresponding redisCluster object, and an error if there is any.
func (c *FakeRedisClusters) Get(name string, options v1.GetOptions) (result *v1alpha1.RedisCluster, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(redisclustersResource, c.ns, name), &v1alpha1.RedisCluster{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.RedisCluster), err
}

// List takes label and field selectors, and returns the list of RedisClusters that match those selectors.
func (c *FakeRedisClusters) List(opts v1.ListOptions) (result *v1alpha1.RedisClusterList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(redisclustersResource, redisclustersKind, c.ns, opts), &v1alpha1.RedisClusterList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.RedisClusterList{}
	for _, item := range obj.(*v1alpha1.RedisClusterList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested redisClusters.
func (c *FakeRedisClusters) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(redisclustersResource, c.ns, opts))

}

// Create takes the representation of a redisCluster and creates it.  Returns the server's representation of the redisCluster, and an error, if there is any.
func (c *FakeRedisClusters) Create(redisCluster *v1alpha1.RedisCluster) (result *v1alpha1.RedisCluster, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(redisclustersResource, c.ns, redisCluster), &v1alpha1.RedisCluster{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.RedisCluster), err
}

// Update takes the representation of a redisCluster and updates it. Returns the server's representation of the redisCluster, and an error, if there is any.
func (c *FakeRedisClusters) Update(redisCluster *v1alpha1.RedisCluster) (result *v1alpha1.RedisCluster, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(redisclustersResource, c.ns, redisCluster), &v1alpha1.RedisCluster{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.RedisCluster), err
}

// Delete takes name of the redisCluster and deletes it. Returns an error if one occurs.
func (c *FakeRedisClusters) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(redisclustersResource, c.ns, name), &v1alpha1.RedisCluster{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeRedisClusters) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(redisclustersResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.RedisClusterList{})
	return err
}

// Patch applies the patch and returns the patched redisCluster.
func (c *FakeRedisClusters) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.RedisCluster, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(redisclustersResource, c.ns, name, data, subresources...), &v1alpha1.RedisCluster{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.RedisCluster), err
}
