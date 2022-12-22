// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1

import (
	quotav1 "github.com/openshift/api/quota/v1"
	v1 "k8s.io/api/core/v1"
)

// ClusterResourceQuotaStatusApplyConfiguration represents an declarative configuration of the ClusterResourceQuotaStatus type for use
// with apply.
type ClusterResourceQuotaStatusApplyConfiguration struct {
	Total      *v1.ResourceQuotaStatus                  `json:"total,omitempty"`
	Namespaces *quotav1.ResourceQuotasStatusByNamespace `json:"namespaces,omitempty"`
}

// ClusterResourceQuotaStatusApplyConfiguration constructs an declarative configuration of the ClusterResourceQuotaStatus type for use with
// apply.
func ClusterResourceQuotaStatus() *ClusterResourceQuotaStatusApplyConfiguration {
	return &ClusterResourceQuotaStatusApplyConfiguration{}
}

// WithTotal sets the Total field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Total field is set to the value of the last call.
func (b *ClusterResourceQuotaStatusApplyConfiguration) WithTotal(value v1.ResourceQuotaStatus) *ClusterResourceQuotaStatusApplyConfiguration {
	b.Total = &value
	return b
}

// WithNamespaces sets the Namespaces field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Namespaces field is set to the value of the last call.
func (b *ClusterResourceQuotaStatusApplyConfiguration) WithNamespaces(value quotav1.ResourceQuotasStatusByNamespace) *ClusterResourceQuotaStatusApplyConfiguration {
	b.Namespaces = &value
	return b
}