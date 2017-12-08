package v1alpha1

import (
	v1alpha1 "github.com/OlivierBoucher/redis-cluster-operator/pkg/apis/rediscluster/v1alpha1"
	"github.com/OlivierBoucher/redis-cluster-operator/pkg/client/clientset/versioned/scheme"
	serializer "k8s.io/apimachinery/pkg/runtime/serializer"
	rest "k8s.io/client-go/rest"
)

type RedisclusterV1alpha1Interface interface {
	RESTClient() rest.Interface
	RedisClustersGetter
}

// RedisclusterV1alpha1Client is used to interact with features provided by the rediscluster.everflow.io group.
type RedisclusterV1alpha1Client struct {
	restClient rest.Interface
}

func (c *RedisclusterV1alpha1Client) RedisClusters(namespace string) RedisClusterInterface {
	return newRedisClusters(c, namespace)
}

// NewForConfig creates a new RedisclusterV1alpha1Client for the given config.
func NewForConfig(c *rest.Config) (*RedisclusterV1alpha1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &RedisclusterV1alpha1Client{client}, nil
}

// NewForConfigOrDie creates a new RedisclusterV1alpha1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *RedisclusterV1alpha1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new RedisclusterV1alpha1Client for the given RESTClient.
func New(c rest.Interface) *RedisclusterV1alpha1Client {
	return &RedisclusterV1alpha1Client{c}
}

func setConfigDefaults(config *rest.Config) error {
	gv := v1alpha1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *RedisclusterV1alpha1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
