package v1alpha1

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +resource:path=rediscluster

type RedisCluster struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object’s metadata
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// Specification of the desired behavior of the Redis cluster
	Spec RedisClusterSpec `json:"spec"`
	// ReadOnly redis information to be used internally
	Status RedisClusterStatus `json:"status"`
}

type RedisClusterSpec struct {
	// Standard object’s metadata
	// Metadata Labels and Annotations gets propagated to the pods.
	PodMetadata *metav1.ObjectMeta `json:"podMetadata,omitempty"`
	// Version tag of the redis image to be used
	Version string `json:"version,omitempty"`
	// Define the number of cluster members (masters)
	Members int32 `json:"members"`
	// Define the number of replicas for each member
	ReplicationFactor *int32 `json:"replicationFactor,omitempty"`
	// Define the volume claim template for each pod
	Storage *v1.PersistentVolumeClaim
	// Define resources requests and limits for single Pods.
	Resources v1.ResourceRequirements `json:"resources,omitempty"`
	// If specified, the pod's scheduling constraints.
	Affinity *v1.Affinity `json:"affinity,omitempty"`
	// If specified, the pod's tolerations.
	Tolerations []v1.Toleration `json:"tolerations,omitempty"`
}

type RedisClusterStatus struct {
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +resource:path=redisclusters

type RedisClusterList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata
	metav1.ListMeta `json:"metadata,omitempty"`
	// List of Redises
	Items []*RedisCluster `json:"items"`
}
