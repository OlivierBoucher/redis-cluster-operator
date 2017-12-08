package fake

import (
	v1alpha1 "github.com/OlivierBoucher/redis-cluster-operator/pkg/client/clientset/versioned/typed/rediscluster/v1alpha1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeRedisclusterV1alpha1 struct {
	*testing.Fake
}

func (c *FakeRedisclusterV1alpha1) RedisClusters(namespace string) v1alpha1.RedisClusterInterface {
	return &FakeRedisClusters{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeRedisclusterV1alpha1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
