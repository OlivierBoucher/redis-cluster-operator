package redis

import (
	"fmt"
	"time"

	clientset "github.com/OlivierBoucher/redis-cluster-operator/pkg/client/clientset/versioned"
	redisscheme "github.com/OlivierBoucher/redis-cluster-operator/pkg/client/clientset/versioned/scheme"
	informers "github.com/OlivierBoucher/redis-cluster-operator/pkg/client/informers/externalversions"
	listers "github.com/OlivierBoucher/redis-cluster-operator/pkg/client/listers/rediscluster/v1alpha1"

	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"

	"go.uber.org/zap"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

const controllerAgentName = "redis-cluster-operator"

type Operator struct {
	kclient kubernetes.Interface
	rclient clientset.Interface

	redisInf    cache.SharedIndexInformer
	redisLister listers.RedisClusterLister
	redisSynced cache.InformerSynced

	queue workqueue.RateLimitingInterface

	// recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	recorder record.EventRecorder

	logger *zap.Logger
}

func NewOperator(kclient *kubernetes.Clientset, rclient *clientset.Clientset) (*Operator, error) {
	logger, err := zap.NewProductionConfig().Build()
	if err != nil {
		return nil, err
	}

	redisInformer := informers.NewSharedInformerFactory(rclient, 5*time.Minute).Rediscluster().V1alpha1().RedisClusters()

	// Create event broadcaster
	// Add redis-cluster-operator types to the default Kubernetes Scheme so Events can be
	// logged for redis-cluster-operator types.
	redisscheme.AddToScheme(scheme.Scheme)
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(logger.Sugar().Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kclient.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	o := &Operator{
		kclient: kclient,
		rclient: rclient,

		redisInf:    redisInformer.Informer(),
		redisLister: redisInformer.Lister(),
		redisSynced: redisInformer.Informer().HasSynced,

		queue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "RedisClusters"),

		recorder: recorder,

		logger: logger,
	}

	o.redisInf.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    o.handleAddRedisCluster,
		UpdateFunc: o.handleUpdateRedisCluster,
		DeleteFunc: o.handleDeleteRedisCluster,
	})

	return o, nil
}

func (o *Operator) Run(stopc <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer o.queue.ShutDown()

	if ok := cache.WaitForCacheSync(stopc, o.redisSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	go o.runWorker()

	go o.redisInf.Run(stopc)

	<-stopc
	return nil
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (o *Operator) runWorker() {
	for o.processNextWorkItem() {
	}
}

func (o *Operator) processNextWorkItem() bool {
	obj, shutdown := o.queue.Get()

	if shutdown {
		return false
	}

	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		// We call Done here so the workqueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the workqueue and attempted again after a back-off
		// period.
		defer o.queue.Done(obj)
		var key string
		var ok bool
		// We expect strings to come off the workqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// workqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// workqueue.
		if key, ok = obj.(string); !ok {
			// As the item in the workqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			o.queue.Forget(obj)
			runtime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		// Run sync, passing it the namespace/name string of the
		// Foo resource to be synced.
		if err := o.sync(key); err != nil {
			return fmt.Errorf("error syncing '%s': %s", key, err.Error())
		}
		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		o.queue.Forget(obj)
		o.logger.Info("Successfully synced", zap.String("key", key))
		return nil
	}(obj)

	if err != nil {
		runtime.HandleError(err)
		return true
	}

	return true
}

func (o *Operator) sync(key string) error {
	// here is where all the heavy lifting happens
	return nil
}

func (o *Operator) handleAddRedisCluster(obj interface{}) {

}

func (o *Operator) handleUpdateRedisCluster(old interface{}, cur interface{}) {

}

func (o *Operator) handleDeleteRedisCluster(obj interface{}) {

}

// enqueueRedisCluster takes a RedisCluster resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than RedisCluster.
func (o *Operator) enqueueRedisCluster(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}
	o.queue.AddRateLimited(key)
}

// handleObject will take any resource implementing metav1.Object and attempt
// to find the RedisCluster resource that 'owns' it. It does this by looking at the
// objects metadata.ownerReferences field for an appropriate OwnerReference.
// It then enqueues that RedisCluster resource to be processed. If the object does not
// have an appropriate OwnerReference, it will simply be skipped.
func (o *Operator) handleObject(obj interface{}) {
	var object metav1.Object
	var ok bool
	if object, ok = obj.(metav1.Object); !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			runtime.HandleError(fmt.Errorf("error decoding object, invalid type"))
			return
		}
		object, ok = tombstone.Obj.(metav1.Object)
		if !ok {
			runtime.HandleError(fmt.Errorf("error decoding object tombstone, invalid type"))
			return
		}
		o.logger.Info("Recovered deleted object from tombstone", zap.String("object", object.GetName()))
	}
	o.logger.Info("Processing object", zap.String("object", object.GetName()))
	if ownerRef := metav1.GetControllerOf(object); ownerRef != nil {
		// If this object is not owned by a RedisCluster, we should not do anything more
		// with it.
		if ownerRef.Kind != "RedisCluster" {
			return
		}

		rc, err := o.redisLister.RedisClusters(object.GetNamespace()).Get(ownerRef.Name)
		if err != nil {
			o.logger.Info("Ignoring orphaned object",
				zap.String("object", object.GetSelfLink()),
				zap.String("owner", ownerRef.Name))
			return
		}

		o.enqueueRedisCluster(rc)
		return
	}
}
