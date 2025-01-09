/*
Copyright 2023.

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

package v1alpha1

import (
	"github.com/opendatahub-io/opendatahub-operator/v2/apis/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	FeatureStoreComponentName = "featurestore"
	// FeatureStoreInstanceName the name of the FeatureStore instance singleton.
	// value should match what's set in the XValidation below
	FeatureStoreInstanceName = "default-" + FeatureStoreComponentName
	FeatureStoreKind         = "FeatureStore"
)

// FeatureStoreCommonSpec spec defines the shared desired state of FeatureStore
type FeatureStoreCommonSpec struct {
	// feature store spec exposed to DSC api
	common.DevFlagsSpec `json:",inline"`

	// Namespace for feature stores to be installed, configurable only once when feature store is enabled, defaults to "odh-feature-stores"
	// +kubebuilder:default="odh-feature-stores"
	// +kubebuilder:validation:Pattern="^([a-z0-9]([-a-z0-9]*[a-z0-9])?)?$"
	// +kubebuilder:validation:MaxLength=63
	StoresNamespace string `json:"storesNamespace,omitempty"`
}

// FeatureStoreSpec defines the desired state of FeatureStore
type FeatureStoreSpec struct {
	// feature store spec exposed to DSC api
	FeatureStoreCommonSpec `json:",inline"`
	//  feature store spec exposed only to internal api
}

// FeatureStoreCommonStatus defines the shared observed state of FeatureStore
type FeatureStoreCommonStatus struct {
	StoresNamespace string `json:"storesNamespace,omitempty"`
}

// FeatureStoreStatus defines the observed state of FeatureStore
type FeatureStoreStatus struct {
	common.Status            `json:",inline"`
	FeatureStoreCommonStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:validation:XValidation:rule="self.metadata.name == 'default-featurestore'",message="FeatureStore name must be default-featurestore"
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type==\"Ready\")].status`,description="Ready"
// +kubebuilder:printcolumn:name="Reason",type=string,JSONPath=`.status.conditions[?(@.type==\"Ready\")].reason`,description="Reason"

// FeatureStore is the Schema for the featurestores API
type FeatureStore struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FeatureStoreSpec   `json:"spec,omitempty"`
	Status FeatureStoreStatus `json:"status,omitempty"`
}

func (c *FeatureStore) GetDevFlags() *common.DevFlags {
	return c.Spec.DevFlags
}

func (c *FeatureStore) GetStatus() *common.Status {
	return &c.Status.Status
}

// +kubebuilder:object:root=true

// FeatureStoreList contains a list of FeatureStore
type FeatureStoreList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FeatureStore `json:"items"`
}

func init() {
	SchemeBuilder.Register(&FeatureStore{}, &FeatureStoreList{})
}

// +kubebuilder:object:generate=true
// +kubebuilder:validation:XValidation:rule="(self.managementState != 'Managed') || (oldSelf.storesNamespace == '') || (oldSelf.managementState != 'Managed')|| (self.storesNamespace == oldSelf.storesNamespace)",message="StoresNamespace is immutable when feature store is Managed"
//nolint:lll

// DSCFeatureStore contains all the configuration exposed in DSC instance for FeatureStore component
type DSCFeatureStore struct {
	// configuration fields common across components
	common.ManagementSpec `json:",inline"`
	// feature store specific field
	FeatureStoreCommonSpec `json:",inline"`
}

// DSCFeatureStoreStatus struct holds the status for the FeatureStore component exposed in the DSC
type DSCFeatureStoreStatus struct {
	common.ManagementSpec     `json:",inline"`
	*FeatureStoreCommonStatus `json:",inline"`
}
